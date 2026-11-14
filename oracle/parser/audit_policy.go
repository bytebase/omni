package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// auditNames is a set of audit option names of one or more words.
type auditNames struct {
	names    map[string]bool
	prefixes map[string]bool // every run of leading words of a name
	maxWords int
}

func newAuditNames(names ...string) *auditNames {
	s := &auditNames{names: map[string]bool{}, prefixes: map[string]bool{}}
	for _, name := range names {
		s.names[name] = true
		words := strings.Fields(name)
		for i := 1; i <= len(words); i++ {
			s.prefixes[strings.Join(words[:i], " ")] = true
		}
		s.maxWords = max(s.maxWords, len(words))
	}
	return s
}

// auditWord returns the unquoted word tok spells, if it is one.
func auditWord(tok Token) (string, bool) {
	if tok.Type == tokIDENT || tok.Type >= kwACCESS {
		return tok.Str, true
	}
	return "", false
}

// matchAuditName returns the longest name of sets that the words from the
// current token on spell, and how many tokens it takes; 0 when none does.
func (p *Parser) matchAuditName(sets ...*auditNames) (string, int) {
	most := 0
	for _, s := range sets {
		most = max(most, s.maxWords)
	}
	toks := append([]Token{p.cur}, p.peekAhead(most-1)...)
	name, n := "", 0
	run := ""
	for i, tok := range toks {
		word, ok := auditWord(tok)
		if !ok {
			break
		}
		if run != "" {
			run += " "
		}
		run += word
		prefix := false
		for _, s := range sets {
			if s.names[run] {
				name, n = run, i+1
			}
			prefix = prefix || s.prefixes[run]
		}
		if !prefix {
			break
		}
	}
	return name, n
}

// parseAuditName consumes the name of sets at the current token, or fails
// there.
func (p *Parser) parseAuditName(sets ...*auditNames) (string, error) {
	name, n := p.matchAuditName(sets...)
	if n == 0 {
		return "", p.syntaxErrorAtCur()
	}
	for range n {
		p.advance()
	}
	return name, nil
}

// isAuditPolicyName reports whether the current token can name an audit
// policy: a quoted identifier or a word that is not reserved.
func (p *Parser) isAuditPolicyName() bool {
	return p.cur.Type == tokQIDENT || p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur)
}

// atStatementEnd reports whether the statement ends at the current token.
func (p *Parser) atStatementEnd() bool {
	return p.cur.Type == tokEOF || p.cur.Type == ';'
}

// parseCreateAuditPolicyStmt parses a CREATE AUDIT POLICY statement. The
// current token follows CREATE AUDIT POLICY.
//
// BNF: oracle/parser/bnf/CREATE-AUDIT-POLICY-Unified-Auditing.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-AUDIT-POLICY-Unified-Auditing.html
//
//	CREATE AUDIT POLICY policy
//	    [ privilege_audit_clause ]
//	    [ action_audit_clause ]...
//	    [ role_audit_clause ]
//	    [ WHEN 'audit_condition' EVALUATE PER { STATEMENT | SESSION | INSTANCE } ]
//	    [ ONLY TOPLEVEL ]
//	    [ CONTAINER = { ALL | CURRENT } ]
//
// Oracle takes the clauses in this order, one ACTIONS of standard actions
// and any number of ACTIONS COMPONENT clauses (parseAuditActionClauses).
// A policy with no audit option at all it refuses after the parse
// (ORA-46373).
func (p *Parser) parseCreateAuditPolicyStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.CreateAuditPolicyStmt{Loc: nodes.Loc{Start: start}}
	if !p.isAuditPolicyName() {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = p.advance().Str
	if p.cur.Type == kwPRIVILEGES {
		privs, err := p.parseAuditPrivileges()
		if err != nil {
			return nil, err
		}
		stmt.Privileges = privs
	}
	if err := p.parseAuditActionClauses(&stmt.Actions, &stmt.ComponentActions); err != nil {
		return nil, err
	}
	if p.isKeywordStr("ROLES") {
		roles, err := p.parseAuditRoles()
		if err != nil {
			return nil, err
		}
		stmt.Roles = roles
	}
	if p.cur.Type == kwWHEN {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.WhenCondition = p.advance().Str
		per, err := p.parseAuditEvaluatePer()
		if err != nil {
			return nil, err
		}
		stmt.EvaluatePer = per
	}
	if p.isKeywordStr("ONLY") {
		p.advance()
		if !p.isKeywordStr("TOPLEVEL") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.OnlyToplevel = true
	}
	if p.isKeywordStr("CONTAINER") {
		v, err := p.parseSecurityContainer()
		if err != nil {
			return nil, err
		}
		stmt.ContainerAll = v
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterAuditPolicyStmt parses an ALTER AUDIT POLICY statement. The
// current token follows ALTER AUDIT POLICY.
//
// BNF: oracle/parser/bnf/ALTER-AUDIT-POLICY-Unified-Auditing.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-AUDIT-POLICY-Unified-Auditing.html
//
//	ALTER AUDIT POLICY policy
//	    { { ADD | DROP } { [ privilege_audit_clause ] [ action_audit_clause ]...
//	                       [ role_audit_clause ]
//	                     | ONLY TOPLEVEL }
//	    | CONDITION { DROP
//	                | 'audit_condition' EVALUATE PER { STATEMENT | SESSION | INSTANCE } }
//	    }
//
// compat.tsv records the statement ending after the policy name, after ADD
// or DROP, and after ADD ONLY as accepted, so omni takes those ends.
func (p *Parser) parseAlterAuditPolicyStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AlterAuditPolicyStmt{Loc: nodes.Loc{Start: start}}
	if !p.isAuditPolicyName() {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = p.advance().Str
	switch {
	case p.atStatementEnd():
	case p.isKeywordStr("ADD") || p.cur.Type == kwDROP:
		stmt.AddDrop = p.advance().Str
		if p.atStatementEnd() {
			break
		}
		if p.isKeywordStr("ONLY") {
			p.advance()
			if p.atStatementEnd() {
				break
			}
			if !p.isKeywordStr("TOPLEVEL") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if stmt.AddDrop == "ADD" {
				stmt.AddToplevel = true
			} else {
				stmt.DropToplevel = true
			}
			break
		}
		found := false
		if p.cur.Type == kwPRIVILEGES {
			privs, err := p.parseAuditPrivileges()
			if err != nil {
				return nil, err
			}
			stmt.Privileges = privs
			found = true
		}
		if p.isKeywordStr("ACTIONS") {
			if err := p.parseAuditActionClauses(&stmt.Actions, &stmt.ComponentActions); err != nil {
				return nil, err
			}
			found = true
		}
		if p.isKeywordStr("ROLES") {
			roles, err := p.parseAuditRoles()
			if err != nil {
				return nil, err
			}
			stmt.Roles = roles
			found = true
		}
		if !found {
			return nil, p.syntaxErrorAtCur()
		}
	case p.isKeywordStr("CONDITION"):
		p.advance()
		switch {
		case p.cur.Type == kwDROP:
			p.advance()
			stmt.ConditionDrop = true
		case p.cur.Type == tokSCONST:
			stmt.Condition = p.advance().Str
			per, err := p.parseAuditEvaluatePer()
			if err != nil {
				return nil, err
			}
			stmt.EvaluatePer = per
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAuditEvaluatePer parses EVALUATE PER { STATEMENT | SESSION |
// INSTANCE } and returns the last word.
func (p *Parser) parseAuditEvaluatePer() (string, error) {
	for _, w := range []string{"EVALUATE", "PER"} {
		if !p.isKeywordStr(w) {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if !p.isKeywordStr("STATEMENT") && !p.isKeywordStr("SESSION") && !p.isKeywordStr("INSTANCE") {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseAuditPrivileges parses PRIVILEGES system_privilege [,
// system_privilege ]...
func (p *Parser) parseAuditPrivileges() ([]string, error) {
	p.advance() // PRIVILEGES
	var privs []string
	for {
		priv, err := p.parseAuditName(auditSystemPrivileges)
		if err != nil {
			return nil, err
		}
		privs = append(privs, priv)
		if p.cur.Type != ',' {
			return privs, nil
		}
		p.advance()
	}
}

// parseAuditRoles parses ROLES role [, role ]... Oracle takes CONNECT and
// RESOURCE as role names, and no other reserved word (ORA-01937).
func (p *Parser) parseAuditRoles() ([]string, error) {
	p.advance() // ROLES
	var roles []string
	for {
		if !p.isAuditRoleName() {
			return nil, p.syntaxErrorAtCur()
		}
		roles = append(roles, p.advance().Str)
		if p.cur.Type != ',' {
			return roles, nil
		}
		p.advance()
	}
}

// isAuditRoleName reports whether the current token can name a role in
// ROLES.
func (p *Parser) isAuditRoleName() bool {
	if p.cur.Type == tokQIDENT || p.cur.Type == kwCONNECT || p.cur.Type == kwRESOURCE {
		return true
	}
	return p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur)
}

// parseAuditActionClauses parses the ACTIONS clauses of an audit policy:
// one of standard actions, and any number of component actions, each led
// by its own ACTIONS.
//
//	action_audit_clause:
//	    ACTIONS { standard_actions | component_actions }
//
//	standard_actions:
//	    { object_action [ ( column [, column ]... ) ] ON object
//	    | ALL ON object
//	    | system_action
//	    | ALL
//	    } [, ... ]...
//
//	object:
//	    { [ schema. ] object_name | DIRECTORY directory_name
//	    | MINING [ MODEL ] [ schema. ] object_name }
func (p *Parser) parseAuditActionClauses(actions *[]*nodes.AuditActionEntry, comps *[]*nodes.AuditComponentAction) error {
	standard := false
	for p.isKeywordStr("ACTIONS") {
		p.advance()
		if p.isKeywordStr("COMPONENT") {
			comp, err := p.parseAuditComponentActions()
			if err != nil {
				return err
			}
			*comps = append(*comps, comp)
			continue
		}
		if standard {
			return p.syntaxErrorAtCur()
		}
		standard = true
		for {
			entry, err := p.parseAuditStandardAction()
			if err != nil {
				return err
			}
			*actions = append(*actions, entry)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	return nil
}

// parseAuditStandardAction parses one standard action.
func (p *Parser) parseAuditStandardAction() (*nodes.AuditActionEntry, error) {
	entry := &nodes.AuditActionEntry{Loc: nodes.Loc{Start: p.pos()}}
	name, err := p.parseAuditName(auditSystemActions, auditObjectActions)
	if err != nil {
		return nil, err
	}
	entry.Action = name
	if auditObjectActions.names[name] {
		if p.cur.Type == '(' {
			p.advance()
			for {
				if !p.isIdentLike() {
					return nil, p.syntaxErrorAtCur()
				}
				entry.Columns = append(entry.Columns, p.advance().Str)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type != ')' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != kwON {
				return nil, p.syntaxErrorAtCur()
			}
		}
		// An object action that is no system action stands with its
		// object.
		if p.cur.Type == kwON || !auditSystemActions.names[name] {
			if p.cur.Type != kwON {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.parseAuditActionObject(entry); err != nil {
				return nil, err
			}
		}
	}
	entry.Loc.End = p.prev.End
	return entry, nil
}

// parseAuditActionObject parses the object after ON in a standard action.
func (p *Parser) parseAuditActionObject(entry *nodes.AuditActionEntry) error {
	switch {
	case p.isKeywordStr("DIRECTORY"):
		p.advance()
		if !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		entry.Directory = p.advance().Str
		return nil
	case p.isKeywordStr("MINING") && p.peekNext().Type != '.' && p.peekNext().Type != tokEOF:
		p.advance()
		if p.cur.Type == kwMODEL {
			p.advance()
		}
		entry.MiningModel = true
	}
	name, err := p.parseAuditObjectName()
	if err != nil {
		return err
	}
	entry.Object = name
	return nil
}

// parseAuditObjectName parses [ schema. ] object_name, with no database
// link: Oracle audits no remote object (ORA-02021).
func (p *Parser) parseAuditObjectName() (*nodes.ObjectName, error) {
	start := p.pos()
	if !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	name := &nodes.ObjectName{Name: p.advance().Str, Loc: nodes.Loc{Start: start}}
	if p.cur.Type == '.' {
		p.advance()
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		name.Schema = name.Name
		name.Name = p.advance().Str
	}
	name.Loc.End = p.prev.End
	return name, nil
}

// auditComponentActions are the actions of the components whose actions
// are not in a dictionary list: Data Pump, direct path loads, the SQL
// Firewall and network protocols.
var (
	auditDatapumpActions    = newAuditNames("ALL", "EXPORT", "IMPORT")
	auditDirectLoadActions  = newAuditNames("ALL", "LOAD")
	auditSQLFirewallActions = newAuditNames("ALL", "CONTEXT VIOLATION", "SQL VIOLATION")
	// FULL is no protocol Oracle audits: it refuses it with ORA-46356, but
	// compat.tsv records PROTOCOL FULL as accepted, so omni takes it.
	auditProtocolActions = newAuditNames("AUTHENTICATION", "FTP", "HTTP", "FULL")
)

// parseAuditComponentActions parses component_actions. The current token is
// COMPONENT.
//
//	COMPONENT = { DATAPUMP { EXPORT | IMPORT | ALL } [, ... ]...
//	            | DIRECT_LOAD { LOAD | ALL } [, ... ]...
//	            | OLS { ols_action | ALL } [, ... ]...
//	            | XS { xs_action | ALL } [, ... ]...
//	            | DV { dv_action ON object_name | ALL } [, ... ]...
//	            | SQL_FIREWALL { ALL | SQL VIOLATION | CONTEXT VIOLATION } ON user [, ... ]...
//	            | PROTOCOL { HTTP | FTP | AUTHENTICATION } [, ... ]... }
//
// The page shows DV's ON optional and SQL_FIREWALL with no action; Oracle
// Free 23.26.3 takes DV actions with or without it, and SQL_FIREWALL
// actions only with it.
func (p *Parser) parseAuditComponentActions() (*nodes.AuditComponentAction, error) {
	comp := &nodes.AuditComponentAction{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // COMPONENT
	if p.cur.Type != '=' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var set *auditNames
	on := 0 // 0: no ON; 1: optional ON; 2: ON required
	switch {
	case p.isKeywordStr("DATAPUMP"):
		set = auditDatapumpActions
	case p.isKeywordStr("DIRECT_LOAD"):
		set = auditDirectLoadActions
	case p.isKeywordStr("OLS"):
		set = auditOLSActions
	case p.isKeywordStr("XS"):
		set = auditXSActions
	case p.isKeywordStr("DV"):
		set, on = auditDVActions, 1
	case p.isKeywordStr("SQL_FIREWALL"):
		set, on = auditSQLFirewallActions, 2
	case p.isKeywordStr("PROTOCOL"):
		set = auditProtocolActions
	default:
		return nil, p.syntaxErrorAtCur()
	}
	comp.Component = p.advance().Str
	for {
		action, err := p.parseAuditName(set)
		if err != nil {
			return nil, err
		}
		comp.Actions = append(comp.Actions, action)
		if on == 2 && p.cur.Type != kwON {
			return nil, p.syntaxErrorAtCur()
		}
		if on > 0 && p.cur.Type == kwON && !(comp.Component == "DV" && action == "ALL") {
			p.advance()
			if p.cur.Type != tokQIDENT && (!p.isIdentLike() || isOracleSQLReservedKeyword(p.cur)) {
				return nil, p.syntaxErrorAtCur()
			}
			comp.Object = p.advance().Str
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	comp.Loc.End = p.prev.End
	return comp, nil
}

// parseAuditPolicyUsers parses what follows BY or EXCEPT in AUDIT and
// NOAUDIT POLICY, and returns the users or, for the second form, the roles:
//
//	{ user [, user ]... | USERS WITH GRANTED { ROLE | ROLES } role [, role ]... }
//
// The page shows USERS WITH ROLE; Oracle Free 23.26.3 refuses ROLE after
// WITH and takes GRANTED ROLES. A user name may be any word: Oracle takes
// AUDIT POLICY p BY SELECT.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/AUDIT-Unified-Auditing.html
func (p *Parser) parseAuditPolicyUsers() (users, roles []string, err error) {
	if p.isKeywordStr("USERS") {
		p.advance()
		for _, w := range []string{"WITH", "GRANTED"} {
			if !p.isKeywordStr(w) {
				return nil, nil, p.syntaxErrorAtCur()
			}
			p.advance()
		}
		if !p.isKeywordStr("ROLE") && !p.isKeywordStr("ROLES") {
			return nil, nil, p.syntaxErrorAtCur()
		}
		p.advance()
		roles, err = p.parseAuditWordList()
		return nil, roles, err
	}
	users, err = p.parseAuditWordList()
	return users, nil, err
}

// parseAuditWordList parses word [, word ]..., the users, roles and
// attributes of AUDIT and NOAUDIT.
func (p *Parser) parseAuditWordList() ([]string, error) {
	var list []string
	for {
		if !p.isIdentLike() {
			return nil, p.syntaxErrorAtCur()
		}
		list = append(list, p.advance().Str)
		if p.cur.Type != ',' {
			return list, nil
		}
		p.advance()
	}
}

// parseAuditContextNamespace parses the context form of AUDIT and NOAUDIT.
// The current token is CONTEXT.
//
//	CONTEXT NAMESPACE namespace ATTRIBUTES attribute [, attribute ]...
//	    [ BY user [, user ]... ]
func (p *Parser) parseAuditContextNamespace() (ns string, attrs, users []string, err error) {
	p.advance() // CONTEXT
	if !p.isKeywordStr("NAMESPACE") {
		return "", nil, nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isIdentLike() {
		return "", nil, nil, p.syntaxErrorAtCur()
	}
	ns = p.advance().Str
	if !p.isKeywordStr("ATTRIBUTES") {
		return "", nil, nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if attrs, err = p.parseAuditWordList(); err != nil {
		return "", nil, nil, err
	}
	if p.cur.Type == kwBY {
		p.advance()
		if users, err = p.parseAuditWordList(); err != nil {
			return "", nil, nil, err
		}
	}
	return ns, attrs, users, nil
}
