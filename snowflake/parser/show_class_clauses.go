package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// SHOW clauses particular to one object class
// ---------------------------------------------------------------------------
//
// The generic SHOW reader (show.go) takes the scopes and filters most pages
// share. The scopes and trailing clauses below belong to one class's page
// each, so they are read only for that class.

// showClassScope is a scope keyword run that a class's page writes after IN.
type showClassScope struct {
	words  []string
	name   showScopeName
	parens bool // the name is written in parentheses
}

// showScopeName says how a class scope's name is written.
type showScopeName int

const (
	scopeNameQualified showScopeName = iota // <name>, possibly qualified (grammar QNAME)
	scopeNameSimple                         // <name>, one part (grammar ID)
	scopeNameOptional                       // [ <name> ], possibly qualified
)

// showClassScopes maps a SHOW class to the scopes only its page writes.
var showClassScopes = map[string][]showClassScope{
	// https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups
	"FAILOVER GROUPS": {{words: []string{"ACCOUNT"}, name: scopeNameSimple}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-replication-groups
	"REPLICATION GROUPS": {{words: []string{"ACCOUNT"}, name: scopeNameSimple}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-roles
	"ROLES": {{words: []string{"CLASS"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-functions
	"FUNCTIONS": {{words: []string{"CLASS"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-procedures
	"PROCEDURES": {{words: []string{"CLASS"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-git-branches
	"GIT BRANCHES": {{words: []string{"GIT", "REPOSITORY"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-git-tags
	"GIT TAGS": {{words: []string{"GIT", "REPOSITORY"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-openflow-connector
	"VERSIONS": {{words: []string{"OPENFLOW", "CONNECTOR"}}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-rules-in-event-routing-table
	"RULES": {{words: []string{"EVENT", "ROUTING", "TABLE"}, parens: true}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-channels
	"CHANNELS": {{words: []string{"PIPE"}, name: scopeNameOptional}},
	// https://docs.snowflake.com/en/sql-reference/sql/show-parameters
	"PARAMETERS": {{words: []string{"PIPE"}, name: scopeNameOptional}},
}

// parseShowClassScope reads, once parseShowScope has consumed IN (or FOR), a
// scope that only stmt's class writes, and reports whether one started at
// cur. A scope whose name is required counts only when a name (or, for a
// parenthesized one, '(') follows its keywords, so a schema of that name
// still reads as one.
func (p *Parser) parseShowClassScope(stmt *ast.ShowStmt) (bool, error) {
	for _, sc := range showClassScopes[stmt.ObjectClass] {
		if !p.atWords(sc.words...) {
			continue
		}
		if stmt.ObjectClass != "PARAMETERS" && p.prev.Type != kwIN {
			continue
		}
		if sc.name != scopeNameOptional && !p.showClassScopeNameFollows(len(sc.words), sc.parens) {
			continue
		}
		p.acceptWords(sc.words...)
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = strings.Join(sc.words, " ")
		var err error
		switch {
		case sc.parens:
			if _, err = p.expect('('); err != nil {
				return true, err
			}
			if stmt.ScopeName, err = p.parseObjectName(); err != nil {
				return true, err
			}
			_, err = p.expect(')')
		case sc.name == scopeNameSimple:
			stmt.ScopeName, err = p.parseObjectSimpleName()
		case sc.name == scopeNameOptional:
			err = p.parseOptionalScopeName(stmt)
		default:
			stmt.ScopeName, err = p.parseObjectName()
		}
		return true, err
	}
	return false, nil
}

// showClassScopeNameFollows reports whether the token after the n keyword
// tokens at cur can start the scope's name ('(' when parens). It consumes
// nothing.
func (p *Parser) showClassScopeNameFollows(n int, parens bool) bool {
	snap := p.snapshot()
	defer p.restore(snap)
	for range n {
		p.advance()
	}
	if parens {
		return p.curIs('(')
	}
	return p.startsScopeName() && !p.curIs(kwLIKE) && !p.curIs(kwFOR)
}

// showClassClause is a trailing clause a class's page writes after its scope.
type showClassClause struct {
	words    []string // the clause's keywords
	scope    string   // the ScopeText the clause follows; "" for none
	bare     bool     // the clause follows IN <name> written without a scope keyword
	required bool     // the page writes the clause unbracketed
	simple   bool     // the name is one part (grammar ID)
}

// showClassClauses maps a SHOW class to its trailing clauses.
var showClassClauses = map[string][]showClassClause{
	// SHOW SEMANTIC DIMENSIONS [ LIKE '<pattern>' ] IN <semantic_view_name> FOR METRIC <metric_name> ...
	// https://docs.snowflake.com/en/sql-reference/sql/show-semantic-dimensions-for-metric
	"SEMANTIC DIMENSIONS": {{words: []string{"FOR", "METRIC"}, bare: true, simple: true}},
	// SHOW VERSIONS [ LIKE '<pattern>' ] IN ARTIFACT REPOSITORY <name> FOR PACKAGE <package_name> [ LIMIT <rows> ]
	// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-artifact-repository
	"VERSIONS": {{words: []string{"FOR", "PACKAGE"}, scope: "ARTIFACT REPOSITORY", required: true, simple: true}},
	// SHOW RELEASE DIRECTIVES [ LIKE '<pattern>' ] IN APPLICATION PACKAGE <name> [ FOR RELEASE CHANNEL <release_channel> ]
	// https://docs.snowflake.com/en/sql-reference/sql/show-release-directives
	"RELEASE DIRECTIVES": {{words: []string{"FOR", "RELEASE", "CHANNEL"}, scope: "APPLICATION PACKAGE", simple: true}},
	// SHOW SHARED CONTENT IN APPLICATION PACKAGE <pkg_name> FOR VERSION <version_name>
	// https://docs.snowflake.com/en/sql-reference/sql/show-shared-content
	"SHARED CONTENT": {{words: []string{"FOR", "VERSION"}, scope: "APPLICATION PACKAGE", required: true, simple: true}},
	// SHOW FUNCTIONS [ LIKE '<pattern>' ] IN MODEL <model_name> [ VERSION <version_name> ]
	// https://docs.snowflake.com/en/sql-reference/sql/show-functions-in-model
	"FUNCTIONS": {{words: []string{"VERSION"}, scope: "MODEL", simple: true}},
	// SHOW RUN { METRICS | PARAMETERS } [ LIKE '<pattern>' ] IN EXPERIMENT <experiment_name> [ RUN <run_name> ] ...
	// https://docs.snowflake.com/en/sql-reference/sql/show-run-in-experiment
	"RUN METRICS":    {{words: []string{"RUN"}, scope: "EXPERIMENT", simple: true}},
	"RUN PARAMETERS": {{words: []string{"RUN"}, scope: "EXPERIMENT", simple: true}},
}

// parseShowClassClause reads the trailing clause stmt's class writes after
// its scope, if any. bareScope reports that the scope was IN <name> written
// without a scope keyword.
func (p *Parser) parseShowClassClause(stmt *ast.ShowStmt, bareScope bool) error {
	for _, cl := range showClassClauses[stmt.ObjectClass] {
		if cl.bare && !bareScope || !cl.bare && (stmt.Scope != ast.ShowScopeOther || stmt.ScopeText != cl.scope) {
			continue
		}
		if !p.acceptWords(cl.words...) {
			if cl.required {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
		stmt.Qualifier = strings.Join(cl.words, " ")
		var err error
		if cl.simple {
			stmt.QualifierName, err = p.parseObjectSimpleName()
		} else {
			stmt.QualifierName, err = p.parseObjectName()
		}
		return err
	}
	if stmt.ObjectClass == "REPLICATION DATABASES" && stmt.Scope == ast.ShowScopeNone && p.atWords("WITH", "PRIMARY") {
		return p.parseShowWithPrimary(stmt)
	}
	return nil
}

// parseShowWithPrimary reads SHOW REPLICATION DATABASES' WITH PRIMARY
// <account_identifier>.<primary_db_name>: the account identifier may itself
// carry the organization (myorg.account1.mydb1), so the name has two or three
// parts.
// https://docs.snowflake.com/en/sql-reference/sql/show-replication-databases
func (p *Parser) parseShowWithPrimary(stmt *ast.ShowStmt) error {
	p.advance() // consume WITH
	p.advance() // consume PRIMARY
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	if name.Schema.IsEmpty() {
		return &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: "syntax error: WITH PRIMARY takes <account_identifier>.<primary_db_name>",
		}
	}
	stmt.Qualifier = "WITH PRIMARY"
	stmt.QualifierName = name
	return nil
}

// parseShowOfType reads SHOW [ JOB ] SERVICES' trailing
// OF TYPE <workload_type> [ , ... ].
// https://docs.snowflake.com/en/sql-reference/sql/show-services
func (p *Parser) parseShowOfType(stmt *ast.ShowStmt) error {
	switch stmt.ObjectClass {
	case "SERVICES", "JOB SERVICES", "SERVICES EXCLUDE JOBS", "JOB SERVICES EXCLUDE JOBS":
	default:
		return nil
	}
	if !p.acceptWords("OF", "TYPE") {
		return nil
	}
	for {
		id, err := p.parseIdent()
		if err != nil {
			return err
		}
		stmt.OfTypes = append(stmt.OfTypes, id)
		if !p.curIs(',') {
			return nil
		}
		p.advance() // consume ','
	}
}

// parseShowGrantsToShareIn reads the IN APPLICATION PACKAGE <app_package_name>
// that SHOW GRANTS TO SHARE <share_name> may carry.
// https://docs.snowflake.com/en/sql-reference/sql/show-grants
func (p *Parser) parseShowGrantsToShareIn(stmt *ast.ShowStmt) error {
	if stmt.GrantsTo == nil || stmt.GrantsTo.Kind != ast.GranteeShare || !p.acceptWords("IN", "APPLICATION", "PACKAGE") {
		return nil
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	stmt.Scope = ast.ShowScopeOther
	stmt.ScopeText = "APPLICATION PACKAGE"
	stmt.ScopeName = name
	return nil
}
