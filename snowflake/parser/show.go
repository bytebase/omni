package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Utility / introspection — SHOW + DESCRIBE / DESC (T6.3)
// ---------------------------------------------------------------------------
//
// Snowflake documents 50+ SHOW object classes and 30+ DESCRIBE object types and
// adds new ones over time. Rather than enumerate them (as the legacy ANTLR
// grammar does with ~68 show_* and ~32 describe_* rules), the object class/type
// is parsed as an open-ended uppercased token run. Structural keywords (TERSE,
// HISTORY, LIKE, IN, FOR, STARTS WITH, LIMIT, WITH, the ->> result pipe, and —
// for DESCRIBE — TYPE) anchor the grammar; the catalog/semantic layer, not the
// parser, validates that a class/type is real. New object classes therefore
// parse without code changes.

// ---------------------------------------------------------------------------
// SHOW
// ---------------------------------------------------------------------------

// parseShowStmt parses a SHOW statement:
//
//	SHOW [TERSE] <object_class> [HISTORY] [LIKE '<pat>']
//	     [IN <scope>] [STARTS WITH '<s>'] [LIMIT <n> [FROM '<s>']]
//	     [WITH PRIVILEGES <priv> [, ...]] [ ->> <query> ]
//	SHOW GRANTS [ON <target> | TO <grantee> | OF <grantee>]
//	SHOW FUTURE GRANTS IN { DATABASE <db> | SCHEMA <schema> }
//
// The SHOW keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseShowStmt() (ast.Node, error) {
	showTok := p.advance() // consume SHOW
	start := showTok.Loc

	stmt := &ast.ShowStmt{Loc: ast.Loc{Start: start.Start}}

	// SHOW GRANTS ... and SHOW FUTURE GRANTS ... are handled specially.
	if p.cur.Type == kwGRANTS {
		p.advance() // consume GRANTS
		if err := p.parseShowGrantsOpts(stmt); err != nil {
			return nil, err
		}
		if err := p.parseShowLimit(stmt); err != nil {
			return nil, err
		}
		if err := p.parseShowPipe(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}
	if p.cur.Type == kwFUTURE && p.peekNext().Type == kwGRANTS {
		p.advance() // consume FUTURE
		p.advance() // consume GRANTS
		stmt.IsGrants = true
		stmt.Future = true
		if err := p.parseShowFutureGrants(stmt); err != nil {
			return nil, err
		}
		if err := p.parseShowLimit(stmt); err != nil {
			return nil, err
		}
		if err := p.parseShowPipe(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Optional TERSE.
	if p.cur.Type == kwTERSE {
		p.advance()
		stmt.Terse = true
	}

	// Object class — an open-ended run of name words.
	class, err := p.parseShowObjectClass()
	if err != nil {
		return nil, err
	}
	stmt.ObjectClass = class

	// Optional HISTORY.
	if p.cur.Type == kwHISTORY {
		p.advance()
		stmt.History = true
	}

	// Optional LIKE '<pat>'.
	if p.cur.Type == kwLIKE {
		p.advance() // consume LIKE
		pat, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		stmt.Like = pat.Str
		stmt.HasLike = true
	}

	// Optional IN/FOR <scope>.
	if p.cur.Type == kwIN || p.cur.Type == kwFOR {
		if err := p.parseShowScope(stmt); err != nil {
			return nil, err
		}
	}

	// SHOW NODES IN COMPUTE POOL <name> [ LIKE '<pattern>' ] writes LIKE after
	// its scope (https://docs.snowflake.com/en/sql-reference/sql/show-nodes-compute-pool).
	if stmt.ObjectClass == "NODES" && !stmt.HasLike && p.cur.Type == kwLIKE {
		p.advance() // consume LIKE
		pat, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		stmt.Like = pat.Str
		stmt.HasLike = true
	}

	// Optional STARTS WITH '<s>'.
	if p.cur.Type == kwSTARTS {
		p.advance() // consume STARTS
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		s, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		stmt.StartsWith = s.Str
		stmt.HasStarts = true
	}

	// Optional LIMIT <n> [FROM '<s>'].
	if p.cur.Type == kwLIMIT {
		p.advance() // consume LIMIT
		n, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		stmt.Limit = n.Str
		stmt.HasLimit = true
		if p.cur.Type == kwFROM {
			p.advance() // consume FROM
			from, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			stmt.LimitFrom = from.Str
		}
	}

	// SHOW AVAILABLE LISTINGS' IS_<x> = TRUE filters (listing.go).
	if stmt.ObjectClass == "AVAILABLE LISTINGS" {
		if err := p.parseShowAvailableListingsFilters(stmt); err != nil {
			return nil, err
		}
	}

	// Optional WITH PRIVILEGES <priv> [, ...].
	if p.cur.Type == kwWITH && p.peekNext().Type == kwPRIVILEGES {
		p.advance() // consume WITH
		p.advance() // consume PRIVILEGES
		privs, err := p.parseShowPrivilegeList()
		if err != nil {
			return nil, err
		}
		stmt.Privileges = privs
	}

	// Optional result pipe: ->> <query>.
	if err := p.parseShowPipe(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	if err := checkShowProjectScope(stmt); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseShowObjectClass parses the open-ended object-class word run after
// SHOW [TERSE]. It collects uppercased name words until a structural keyword
// (HISTORY, LIKE, IN, FOR, STARTS, LIMIT, WITH, the ->> pipe) or EOF. Requires
// at least one word.
func (p *Parser) parseShowObjectClass() (string, error) {
	if !p.startsShowClassWord() {
		return "", p.syntaxErrorAtCur()
	}
	var words []string
	for p.startsShowClassWord() {
		tok := p.advance()
		words = append(words, strings.ToUpper(tok.Str))
	}
	return strings.Join(words, " "), nil
}

// startsShowClassWord reports whether the current token continues the object
// class run — i.e. it is a name word and not one of the structural keywords that
// terminate the class.
func (p *Parser) startsShowClassWord() bool {
	switch p.cur.Type {
	case kwHISTORY, kwLIKE, kwIN, kwFOR, kwSTARTS, kwLIMIT, kwWITH, tokFlow, tokEOF:
		return false
	}
	if p.peekNext().Type == '=' {
		// A filter's name (SHOW AVAILABLE LISTINGS IS_IMPORTED = TRUE).
		return false
	}
	return p.isObjectTypeWord(p.cur.Type)
}

// parseShowScope parses the IN/FOR scope qualifier of a generic SHOW:
//
//	IN ACCOUNT
//	IN DATABASE [<name>]
//	IN SCHEMA [<name>]
//	IN TABLE [<name>] | IN VIEW [<name>]
//	{IN|FOR} { SESSION | USER [<name>] | WAREHOUSE [<name>] | TASK [<name>]
//	          | APPLICATION [<name>] | CONNECTION [<name>] }   (ShowScopeOther)
//	IN <schema-name>                 (bare name → schema scope)
//
// The IN/FOR keyword is the current token on entry. The SESSION/USER/WAREHOUSE/
// TASK/... scopes come from SHOW PARAMETERS and friends; they are captured as
// ShowScopeOther with the scope keyword in ScopeText (rather than mis-modeled as
// a bare schema name).
func (p *Parser) parseShowScope(stmt *ast.ShowStmt) error {
	p.advance() // consume IN/FOR

	// SHOW BACKUPS IN BACKUP SET <name> and SHOW SNAPSHOTS IN SNAPSHOT SET
	// <name> (https://docs.snowflake.com/en/sql-reference/sql/show-backups-in-backup-set,
	// https://docs.snowflake.com/en/sql-reference/sql/show-snapshots-in-snapshot-set).
	if (stmt.ObjectClass == "BACKUPS" && p.atWords("BACKUP", "SET")) ||
		(stmt.ObjectClass == "SNAPSHOTS" && p.atWords("SNAPSHOT", "SET")) {
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = strings.ToUpper(p.advance().Str) + " SET"
		p.advance() // consume SET
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		stmt.ScopeName = name
		return nil
	}

	// IN { MODEL | DATASET | EXPERIMENT | GATEWAY } <name> (model.go).
	if word := p.startsShowModelScope(); word != "" {
		return p.parseShowModelScope(stmt, word)
	}

	switch p.cur.Type {
	case kwACCOUNT:
		p.advance()
		stmt.Scope = ast.ShowScopeAccount
		return nil
	case kwDATABASE:
		p.advance()
		stmt.Scope = ast.ShowScopeDatabase
		return p.parseOptionalScopeName(stmt)
	case kwSCHEMA:
		p.advance()
		stmt.Scope = ast.ShowScopeSchema
		return p.parseOptionalScopeName(stmt)
	case kwTABLE:
		p.advance()
		stmt.Scope = ast.ShowScopeTable
		return p.parseOptionalScopeName(stmt)
	case kwVIEW:
		p.advance()
		stmt.Scope = ast.ShowScopeView
		return p.parseOptionalScopeName(stmt)
	case kwEXTERNAL:
		// FOR EXTERNAL CONSUMER <name> (sharing.go); otherwise a schema
		// named EXTERNAL.
		if isWordToken(p.peekNext(), "CONSUMER") {
			return p.parseShowExternalConsumerScope(stmt)
		}
	case kwSESSION, kwUSER, kwWAREHOUSE, kwTASK, kwAPPLICATION, kwCONNECTION:
		// Non-container scope keyword (SHOW PARAMETERS / SHOW ... contexts).
		tok := p.advance()
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = strings.ToUpper(tok.Str)
		// IN APPLICATION PACKAGE <name> (SHOW VERSIONS, SHOW RELEASE
		// DIRECTIVES, SHOW RELEASE CHANNELS), where PACKAGE is followed by
		// the package's name rather than being the application's.
		// https://docs.snowflake.com/en/sql-reference/sql/show-versions
		if tok.Type == kwAPPLICATION && p.startsShowApplicationPackageScope() {
			p.advance() // consume PACKAGE
			stmt.ScopeText = "APPLICATION PACKAGE"
		}
		// SESSION takes no name; the others take an optional name.
		if tok.Type != kwSESSION {
			return p.parseOptionalScopeName(stmt)
		}
		return nil
	case kwFAILOVER, kwREPLICATION:
		// IN FAILOVER GROUP <name> | IN REPLICATION GROUP <name>
		// (SHOW DATABASES / SHOW SHARES scopes). Without this case the
		// bare-name fallback below consumed FAILOVER/REPLICATION as a schema
		// name and left `GROUP <name>` behind as silently-dropped trailing
		// tokens.
		tok := p.advance()
		if _, err := p.expect(kwGROUP); err != nil {
			return err
		}
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = strings.ToUpper(tok.Str) + " GROUP"
		return p.parseOptionalScopeName(stmt)
	}

	// IN { DBT | DCM } PROJECT <name> (project.go).
	if words := p.startsShowProjectScope(stmt.ObjectClass); words != nil {
		return p.parseShowProjectScope(stmt, words)
	}

	// SHOW COMMITTED TASKS ... IN GRAPH <graph_name>
	// https://docs.snowflake.com/en/sql-reference/sql/show-committed-tasks
	// (IN GRAPH alone, or before STARTS WITH or LIMIT, is the schema GRAPH.)
	if next := p.peekNext(); stmt.ObjectClass == "COMMITTED TASKS" && isWordToken(p.cur, "GRAPH") &&
		(isNameWordTok(next) || next.Type == tokQuotedIdent) &&
		next.Type != kwSTARTS && next.Type != kwLIMIT && next.Type != kwWITH {
		p.advance() // consume GRAPH
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = "GRAPH"
		stmt.ScopeName = name
		return nil
	}

	// IN LISTING <name> (SHOW RELEASE CHANNELS); a bare LISTING is a schema.
	// https://docs.snowflake.com/en/sql-reference/sql/show-release-channels
	if p.startsShowListingScope() {
		p.advance() // consume LISTING
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = "LISTING"
		return p.parseOptionalScopeName(stmt)
	}

	// IN SERVICE | COMPUTE POOL | IMAGE REPOSITORY | ARTIFACT REPOSITORY <name>
	// (SHOW ENDPOINTS / ROLES / SERVICE CONTAINERS IN SERVICE, SHOW NODES /
	// SERVICES IN COMPUTE POOL, SHOW IMAGES, SHOW FILES / PACKAGES).
	if words := p.startsShowServiceScope(); words != nil {
		p.acceptWords(words...)
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = strings.Join(words, " ")
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		stmt.ScopeName = name
		return nil
	}

	// Bare name → schema scope (e.g. SHOW TABLES IN tpch_sf1).
	if p.startsScopeName() {
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		stmt.Scope = ast.ShowScopeSchema
		stmt.ScopeName = name
		return nil
	}

	return p.syntaxErrorAtCur()
}

// startsShowServiceScope reports the Snowpark Container Services scope
// keywords at cur (after IN), or nil. SERVICE counts only when a name word
// follows, so SHOW ... IN service still reads a schema named SERVICE. It
// consumes nothing.
func (p *Parser) startsShowServiceScope() []string {
	for _, words := range [][]string{{"COMPUTE", "POOL"}, {"IMAGE", "REPOSITORY"}, {"ARTIFACT", "REPOSITORY"}} {
		if p.atWords(words...) {
			return words
		}
	}
	if p.curIsWord("SERVICE") {
		switch next := p.peekNext().Type; next {
		case kwSTARTS, kwLIMIT, kwWITH, kwLIKE, kwFOR, tokFlow, tokEOF, ';':
		default:
			if p.isObjectTypeWord(next) {
				return []string{"SERVICE"}
			}
		}
	}
	return nil
}

// parseOptionalScopeName reads an optional object name following a scope keyword
// (DATABASE/SCHEMA/TABLE/VIEW), storing it in stmt.ScopeName if present.
func (p *Parser) parseOptionalScopeName(stmt *ast.ShowStmt) error {
	if p.startsScopeName() {
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		stmt.ScopeName = name
	}
	return nil
}

// startsScopeName reports whether the current token can begin a scope object
// name. It excludes the structural keywords that may follow a scope clause
// (STARTS, LIMIT, the pipe, EOF) so an absent name is handled gracefully.
func (p *Parser) startsScopeName() bool {
	switch p.cur.Type {
	case kwSTARTS, kwLIMIT, kwWITH, tokFlow, tokEOF, ';':
		return false
	}
	return p.isObjectTypeWord(p.cur.Type)
}

// parseShowPrivilegeList parses the WITH PRIVILEGES <priv> [, ...] list, reusing
// the GRANT privilege grammar (each privilege is an open-ended uppercased word
// run). The list ends at the pipe or EOF.
func (p *Parser) parseShowPrivilegeList() ([]*ast.Privilege, error) {
	var privs []*ast.Privilege
	for {
		priv, err := p.parseShowPrivilege()
		if err != nil {
			return nil, err
		}
		privs = append(privs, priv)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return privs, nil
}

// parseShowPrivilege parses one privilege word run in a WITH PRIVILEGES list. It
// mirrors parsePrivilege but terminates on a comma or the result pipe rather
// than on ON.
func (p *Parser) parseShowPrivilege() (*ast.Privilege, error) {
	if !p.isPrivilegeWord(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	startLoc := p.cur.Loc
	var b strings.Builder
	first := p.advance()
	b.WriteString(strings.ToUpper(first.Str))
	endLoc := first.Loc
	for {
		if p.cur.Type == ',' || p.cur.Type == tokFlow || p.cur.Type == tokEOF {
			break
		}
		if !p.isPrivilegeWord(p.cur.Type) {
			break
		}
		part := p.advance()
		b.WriteByte(' ')
		b.WriteString(strings.ToUpper(part.Str))
		endLoc = part.Loc
	}
	return &ast.Privilege{
		Name: b.String(),
		Loc:  ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}, nil
}

// parseShowPipe parses an optional trailing result pipe: ->> <query>. The piped
// statement is parsed as a full top-level statement (typically a SELECT).
func (p *Parser) parseShowPipe(stmt *ast.ShowStmt) error {
	if p.cur.Type != tokFlow {
		return nil
	}
	p.advance() // consume ->>
	node, err := p.parseStmt()
	if err != nil {
		return err
	}
	stmt.Pipe = node
	return nil
}

// ---------------------------------------------------------------------------
// SHOW GRANTS
// ---------------------------------------------------------------------------

// parseShowGrantsOpts parses the optional tail of SHOW GRANTS:
//
//	ON ACCOUNT
//	ON <object_type> <name>
//	TO { ROLE | USER | SHARE | DATABASE ROLE | APPLICATION [ROLE] | <class> ROLE } <name>
//	OF { ROLE | SHARE } <name>
//
// On entry GRANTS has already been consumed. stmt.IsGrants is set here.
func (p *Parser) parseShowGrantsOpts(stmt *ast.ShowStmt) error {
	stmt.IsGrants = true

	switch p.cur.Type {
	case kwON:
		p.advance() // consume ON
		target, err := p.parseGrantTarget()
		if err != nil {
			return err
		}
		stmt.GrantsOn = target
		return nil
	case kwTO:
		p.advance() // consume TO
		grantee, err := p.parseGrantee()
		if err != nil {
			return err
		}
		stmt.GrantsTo = grantee
		return nil
	case kwOF:
		p.advance() // consume OF
		grantee, err := p.parseGrantee()
		if err != nil {
			return err
		}
		stmt.GrantsTo = grantee
		return nil
	case kwIN:
		// IN DCM PROJECT <name> (project.go).
		p.advance() // consume IN
		words := p.startsShowProjectScope("GRANTS")
		if words == nil {
			return p.syntaxErrorAtCur()
		}
		return p.parseShowProjectScope(stmt, words)
	}
	// Bare SHOW GRANTS — no options.
	return nil
}

// parseShowFutureGrants parses the tail of SHOW FUTURE GRANTS:
//
//	IN { DATABASE <db> | SCHEMA <schema> }
//	TO { ROLE <name> | DATABASE ROLE <name> }
//
// On entry FUTURE GRANTS has been consumed. The IN container is modeled as a
// GrantTarget of kind GrantTargetAllIn (the future-grants scope is a
// container); the TO grantee reuses the GRANT/REVOKE grantee grammar. Documented
// by Snowflake (the legacy grammar only had the IN forms).
func (p *Parser) parseShowFutureGrants(stmt *ast.ShowStmt) error {
	switch p.cur.Type {
	case kwIN:
		p.advance() // consume IN
		// IN DCM PROJECT <name> (project.go).
		if words := p.startsShowProjectScope("GRANTS"); words != nil {
			return p.parseShowProjectScope(stmt, words)
		}
		target := &ast.GrantTarget{Kind: ast.GrantTargetAllIn}
		switch p.cur.Type {
		case kwDATABASE:
			p.advance() // consume DATABASE
			target.Container = ast.GrantContainerDatabase
		case kwSCHEMA:
			p.advance() // consume SCHEMA
			target.Container = ast.GrantContainerSchema
		default:
			return p.syntaxErrorAtCur()
		}
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		target.ContainerName = name
		target.Loc = name.Loc
		stmt.GrantsOn = target
		return nil

	case kwTO:
		p.advance() // consume TO
		grantee, err := p.parseGrantee()
		if err != nil {
			return err
		}
		stmt.GrantsTo = grantee
		return nil

	default:
		return p.syntaxErrorAtCur()
	}
}

// parseShowLimit parses an optional trailing LIMIT <rows> clause on a SHOW
// GRANTS / SHOW FUTURE GRANTS statement. (Documented by Snowflake; absent from
// the legacy grammar's show_grants rule.) Reuses the ShowStmt.Limit fields.
func (p *Parser) parseShowLimit(stmt *ast.ShowStmt) error {
	if p.cur.Type != kwLIMIT {
		return nil
	}
	p.advance() // consume LIMIT
	n, err := p.expect(tokInt)
	if err != nil {
		return err
	}
	stmt.Limit = n.Str
	stmt.HasLimit = true
	return nil
}

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

// parseDescribeStmt parses a DESCRIBE or DESC statement:
//
//	{ DESCRIBE | DESC } <object_type> <name> [ ( <signature> ) ] [ TYPE = <kw> ]
//	{ DESCRIBE | DESC } SEARCH OPTIMIZATION ON <name>
//	{ DESCRIBE | DESC } RESULT { '<query_id>' | LAST_QUERY_ID() }
//
// short reports whether the keyword was DESC (true) or DESCRIBE (false). The
// keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseDescribeStmt(short bool) (ast.Node, error) {
	descTok := p.advance() // consume DESCRIBE/DESC
	start := descTok.Loc

	stmt := &ast.DescribeStmt{Short: short, Loc: ast.Loc{Start: start.Start}}

	// Object type — an open-ended run of name "units" (dotted identifier paths).
	// Mirroring grant_revoke's parseGrantTargetObject: read units left-to-right,
	// shifting each prior unit into the type, so the LAST unit is the object name
	// and every preceding unit forms the (multi-word) type. The run stops at a
	// terminator: '(' (a signature), TYPE, ON (SEARCH OPTIMIZATION ON ...), a
	// literal (RESULT '<id>' / TRANSACTION <n>), or EOF.
	if !p.isObjectTypeWord(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	var typeWords []string
	unit, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	for p.startsDescribeNameUnit() {
		typeWords = append(typeWords, strings.ToUpper(unit.String()))
		unit, err = p.parseObjectName()
		if err != nil {
			return nil, err
		}
	}

	// Decide whether the final `unit` is the object name or the tail of the type.
	switch p.cur.Type {
	case kwON:
		// SEARCH OPTIMIZATION ON <name>: the whole unit run is the type; the name
		// is introduced by ON.
		typeWords = append(typeWords, strings.ToUpper(unit.String()))
		stmt.ObjectType = strings.Join(typeWords, " ")
		p.advance() // consume ON
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil

	case tokString, tokInt, tokFloat:
		// RESULT '<id>' / TRANSACTION <n>: the whole unit run is the type; the
		// name is a literal.
		typeWords = append(typeWords, strings.ToUpper(unit.String()))
		stmt.ObjectType = strings.Join(typeWords, " ")
		lit, err := p.parseDescribeNameLiteral()
		if err != nil {
			return nil, err
		}
		stmt.NameLiteral = lit
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Standard form: the final unit is the object name.
	if len(typeWords) == 0 {
		return nil, &ParseError{
			Position: start.Start, End: start.End,
			Message: "expected object type before object name in DESCRIBE",
		}
	}
	stmt.ObjectType = strings.Join(typeWords, " ")
	stmt.Name = unit

	// Optional FUNCTION/PROCEDURE argument-type signature.
	if p.cur.Type == '(' {
		sig, err := p.parseGrantSignature()
		if err != nil {
			return nil, err
		}
		stmt.Signature = sig
	}

	// DESC LISTING <name> REVISION = { DRAFT | PUBLISHED } (listing.go).
	if stmt.ObjectType == "LISTING" && p.startsListingRevision() {
		if err := p.parseDescribeListingRevision(stmt); err != nil {
			return nil, err
		}
	}

	// Optional TYPE = <kw>.
	if p.cur.Type == kwTYPE {
		p.advance() // consume TYPE
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if !p.isObjectTypeWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		opt := p.advance()
		stmt.TypeOption = strings.ToUpper(opt.Str)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsDescribeNameUnit reports whether the current token can begin another
// space-separated name unit inside the DESCRIBE type/name run — i.e. it is a
// name word and not a terminator. Terminators: '(' (signature), TYPE, ON
// (SEARCH OPTIMIZATION ON), a literal (RESULT/TRANSACTION), and EOF.
func (p *Parser) startsDescribeNameUnit() bool {
	switch p.cur.Type {
	case '(', kwTYPE, kwON, tokString, tokInt, tokFloat, tokEOF, ';':
		return false
	}
	if p.startsListingRevision() {
		return false
	}
	return p.isObjectTypeWord(p.cur.Type)
}

// parseDescribeNameLiteral parses a literal name for DESCRIBE RESULT '<id>' or
// DESCRIBE TRANSACTION <num>.
func (p *Parser) parseDescribeNameLiteral() (*ast.Literal, error) {
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}, nil
	case tokInt:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}, nil
	case tokFloat:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}, nil
	}
	return nil, p.syntaxErrorAtCur()
}
