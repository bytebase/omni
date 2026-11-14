package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Policies read by their docs pages: AGGREGATION, JOIN, PROJECTION, PRIVACY,
// PACKAGES, FEATURE, MAINTENANCE and INGRESS POLICY
// ---------------------------------------------------------------------------
//
// These policies share ast.CreatePolicyStmt / ast.AlterPolicyStmt with the
// MASKING / ROW ACCESS / SESSION / ... policies of role_user_policy.go, but
// each is read strictly by the form its page gives: a property its page does
// not document, a value of another shape, or a clause the page writes for
// another policy is a syntax error.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-aggregation-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-join-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-projection-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-privacy-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-packages-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-feature-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-maintenance-policy
// https://docs.snowflake.com/en/sql-reference/sql/create-ingress-policy
// and the ALTER and DROP page of each.

// docPolicy is what one policy's pages document.
type docPolicy struct {
	word string
	kind ast.PolicyKind
	// returns is the constraint type a body policy's RETURNS names; "" for
	// a policy of properties.
	returns string
	// create and set are the properties CREATE and ALTER ... SET take.
	create, set policyProps
	// unset are the properties ALTER ... UNSET takes; unsetList says UNSET
	// takes a comma list of them rather than one.
	unset     []string
	unsetList bool
	// rename, tags and definition say the policy documents RENAME TO,
	// SET TAG / UNSET TAG and an AS '<definition>' body.
	rename, tags, definition bool
	// tagsIfExists says the page writes IF EXISTS on SET TAG / UNSET TAG.
	tagsIfExists bool
	// dropIfExists says DROP's page writes IF EXISTS.
	dropIfExists bool
}

// policyProps maps a documented property's name to the values it takes.
type policyProps map[string]func(*ast.CopyOption) bool

func polString(o *ast.CopyOption) bool { return propString.matches(o) }

// polStrings takes ( '<string>' [ , ... ] ), or ( ) when the page writes the
// list as optional.
func polStrings(empty bool) func(*ast.CopyOption) bool {
	return func(o *ast.CopyOption) bool {
		return (empty && isEmptyGroup(o)) || propStrings.matches(o)
	}
}

// featurePolicyObjectTypes are the object types CREATE FEATURE POLICY lists
// for BLOCKED_OBJECT_TYPES_FOR_CREATION, each "accepted in either singular or
// plural form".
// https://docs.snowflake.com/en/sql-reference/sql/create-feature-policy
var featurePolicyObjectTypes = func() map[string]bool {
	m := map[string]bool{}
	for _, pair := range [][2]string{
		{"AGENT", "AGENTS"},
		{"APPLICATION_SERVICE", "APPLICATION_SERVICES"},
		{"ARTIFACT_REPOSITORY", "ARTIFACT_REPOSITORIES"},
		{"COMPUTE_POOL", "COMPUTE_POOLS"},
		{"CORTEX_SEARCH_SERVICE", "CORTEX_SEARCH_SERVICES"},
		{"DATABASE", "DATABASES"},
		{"DYNAMIC_TABLE", "DYNAMIC_TABLES"},
		{"FILE_FORMAT", "FILE_FORMATS"},
		{"FUNCTION", "FUNCTIONS"},
		{"GIT_REPOSITORY", "GIT_REPOSITORIES"},
		{"ICEBERG_TABLE", "ICEBERG_TABLES"},
		{"MATERIALIZED_VIEW", "MATERIALIZED_VIEWS"},
		{"MCP_SERVER", "MCP_SERVERS"},
		{"PROCEDURE", "PROCEDURES"},
		{"SCHEMA", "SCHEMAS"},
		{"SECRET", "SECRETS"},
		{"STAGE", "STAGES"},
		{"TABLE", "TABLES"},
		{"TASK", "TASKS"},
		{"VIEW", "VIEWS"},
		{"WAREHOUSE", "WAREHOUSES"},
		{"WORKSPACE", "WORKSPACES"},
	} {
		m[pair[0]], m[pair[1]] = true, true
	}
	return m
}()

// polObjectTypes takes ( <object_type> [ , ... ] ) of the unquoted object
// types featurePolicyObjectTypes names, or ( ) when empty allows it.
func polObjectTypes(empty bool) func(*ast.CopyOption) bool {
	return func(o *ast.CopyOption) bool {
		if empty && isEmptyGroup(o) {
			return true
		}
		if len(o.Group) == 0 {
			return false
		}
		for _, e := range o.Group {
			if !e.Bare || !featurePolicyObjectTypes[e.Name] {
				return false
			}
		}
		return true
	}
}

func isEmptyGroup(o *ast.CopyOption) bool {
	return o.Group != nil && len(o.Group) == 0 && o.List == nil
}

var commentOnly = policyProps{"COMMENT": polString}

// bodyPolicy is a policy whose CREATE is AS ( ) RETURNS <type> -> <body>.
func bodyPolicy(word string, kind ast.PolicyKind, returns string) *docPolicy {
	return &docPolicy{
		word: word, kind: kind, returns: returns,
		create: commentOnly, set: commentOnly, unset: []string{"COMMENT"},
		rename: true, tags: true,
	}
}

var packagesPolicyProps = policyProps{
	"ALLOWLIST":                     polStrings(true),
	"BLOCKLIST":                     polStrings(true),
	"ADDITIONAL_CREATION_BLOCKLIST": polStrings(true),
	"COMMENT":                       polString,
}

var docPolicies = []*docPolicy{
	bodyPolicy("AGGREGATION", ast.PolicyAggregation, "AGGREGATION_CONSTRAINT"),
	bodyPolicy("JOIN", ast.PolicyJoin, "JOIN_CONSTRAINT"),
	bodyPolicy("PROJECTION", ast.PolicyProjection, "PROJECTION_CONSTRAINT"),
	bodyPolicy("PRIVACY", ast.PolicyPrivacy, "PRIVACY_BUDGET"),
	{
		word: "PACKAGES", kind: ast.PolicyPackages,
		create: packagesPolicyProps, set: packagesPolicyProps,
		unset: []string{"ALLOWLIST", "BLOCKLIST", "ADDITIONAL_CREATION_BLOCKLIST", "COMMENT"}, unsetList: true,
		dropIfExists: true,
	},
	{
		word: "FEATURE", kind: ast.PolicyFeature,
		create: policyProps{"BLOCKED_OBJECT_TYPES_FOR_CREATION": polObjectTypes(false), "COMMENT": polString},
		set:    policyProps{"BLOCKED_OBJECT_TYPES_FOR_CREATION": polObjectTypes(true), "COMMENT": polString},
		unset:  []string{"BLOCKED_OBJECT_TYPES_FOR_CREATION", "COMMENT", "POLICY_DEFINITION"}, unsetList: true,
		rename: true, tags: true, tagsIfExists: true, definition: true,
	},
	{
		word: "MAINTENANCE", kind: ast.PolicyMaintenance,
		create: policyProps{"SCHEDULE": polString, "COMMENT": polString},
		set:    policyProps{"SCHEDULE": polString, "COMMENT": polString},
		unset:  []string{"COMMENT"},
		rename: true, dropIfExists: true,
	},
	{
		word: "INGRESS", kind: ast.PolicyIngress,
		create: commentOnly, set: commentOnly, unset: []string{"COMMENT"},
		definition: true, dropIfExists: true,
	},
}

// startsDocPolicy reports the policy whose <kind> POLICY keywords start at
// cur, when it is one this file reads. It consumes nothing.
func (p *Parser) startsDocPolicy() *docPolicy {
	if p.cur.Type != tokIdent && p.cur.Type < 700 {
		return nil
	}
	for _, d := range docPolicies {
		if p.atWords(d.word, "POLICY") {
			return d
		}
	}
	return nil
}

// parseCreateDocPolicy reads CREATE of policy d. cur is d's kind word; start
// is the CREATE token's Loc.
//
//	CREATE [ OR REPLACE ] { AGGREGATION | JOIN | PROJECTION | PRIVACY } POLICY
//	  [ IF NOT EXISTS ] <name> AS ( ) RETURNS <type> -> <body>
//	  [ COMMENT = '<string>' ]
//	CREATE [ OR REPLACE ] PACKAGES POLICY [ IF NOT EXISTS ] <name>
//	  LANGUAGE PYTHON [ ALLOWLIST = ( [ '<pkg>' [ , ... ] ] ) ]
//	  [ BLOCKLIST = ( ... ) ] [ ADDITIONAL_CREATION_BLOCKLIST = ( ... ) ]
//	  [ COMMENT = '<string>' ]
//	CREATE [ OR REPLACE ] FEATURE POLICY [ IF NOT EXISTS ] <name>
//	  [ BLOCKED_OBJECT_TYPES_FOR_CREATION = ( <type> [ , ... ] ) ]
//	  [ COMMENT = '<string>' ] [ AS '<definition>' ]
//	CREATE [ OR REPLACE ] MAINTENANCE POLICY [ IF NOT EXISTS ] <name>
//	  SCHEDULE = '<cron>' [ COMMENT = '<string>' ]
//	CREATE [ OR REPLACE ] INGRESS POLICY [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string>' ] AS '<definition>'
//
// IF NOT EXISTS excludes OR REPLACE, as each page writes the two forms apart.
func (p *Parser) parseCreateDocPolicy(start ast.Loc, orReplace, orAlter bool, d *docPolicy) (ast.Node, error) {
	if orAlter {
		return p.unsupported("CREATE")
	}
	p.advance() // consume the kind word
	p.advance() // consume POLICY
	stmt := &ast.CreatePolicyStmt{Kind: d.kind, OrReplace: orReplace, Loc: ast.Loc{Start: start.Start}}
	if !orReplace {
		if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
			return nil, err
		}
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if d.returns != "" {
		// AS ( ) RETURNS <type> -> <body>: the signature is always empty.
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwRETURNS); err != nil {
			return nil, err
		}
		if !isWordToken(p.cur, d.returns) {
			return nil, p.syntaxErrorAtCur()
		}
		ret := p.advance()
		stmt.Returns = &ast.TypeName{Kind: ast.TypeUnknown, Name: d.returns, VectorDim: -1, Loc: ret.Loc}
		if _, err := p.expect(tokArrow); err != nil {
			return nil, err
		}
		body, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Body = body
	}
	if d.kind == ast.PolicyPackages {
		if err := p.expectWords("LANGUAGE", "PYTHON"); err != nil {
			return nil, err
		}
		stmt.Language = "PYTHON"
	}

	opts, err := p.parsePolicyProps(d.create, d.kind == ast.PolicyMaintenance)
	if err != nil {
		return nil, err
	}
	stmt.Options = opts
	if d.kind == ast.PolicyMaintenance && !hasOption(opts, "SCHEDULE") {
		return nil, &ParseError{Position: p.cur.Loc.Start, End: p.cur.Loc.End, Message: "syntax error: SCHEDULE is required"}
	}

	if d.definition && (p.cur.Type == kwAS || d.kind == ast.PolicyIngress) {
		def, err := p.parsePolicyDefinition()
		if err != nil {
			return nil, err
		}
		stmt.Definition = def
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterDocPolicy reads ALTER of policy d. cur is d's kind word; start is
// the ALTER token's Loc.
//
//	ALTER <kind> POLICY [ IF EXISTS ] <name>
//	  { RENAME TO <new_name>                          -- not PACKAGES, INGRESS
//	  | SET BODY -> <body>                            -- body policies
//	  | SET <property> = <value> [ ... ] [ AS '<definition>' ]
//	  | AS '<definition>'                             -- FEATURE, INGRESS
//	  | UNSET <property> [ , ... ]
//	  | SET TAG <tag> = '<value>' [ , ... ]           -- body policies, FEATURE
//	  | UNSET TAG <tag> [ , ... ] }
//
// The body policies' pages write no IF EXISTS on SET TAG and UNSET TAG.
func (p *Parser) parseAlterDocPolicy(start ast.Loc, d *docPolicy) (ast.Node, error) {
	p.advance() // consume the kind word
	p.advance() // consume POLICY
	stmt := &ast.AlterPolicyStmt{Kind: d.kind, Loc: ast.Loc{Start: start.Start}}
	if err := p.parseIfExists(&stmt.IfExists); err != nil {
		return nil, err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	tagStart := p.cur
	switch {
	case p.cur.Type == kwRENAME && d.rename:
		newName, err := p.parseRenameTo()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicyRename
		stmt.NewName = newName

	case p.cur.Type == kwSET && p.peekNext().Type == kwTAG && d.tags:
		if stmt.IfExists && !d.tagsIfExists {
			return nil, syntaxErrorAtTok(tagStart)
		}
		p.advance() // consume SET
		p.advance() // consume TAG
		tags, err := p.parseObjectTagList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicySetTag
		stmt.Tags = tags

	case p.cur.Type == kwUNSET && p.peekNext().Type == kwTAG && d.tags:
		if stmt.IfExists && !d.tagsIfExists {
			return nil, syntaxErrorAtTok(tagStart)
		}
		p.advance() // consume UNSET
		p.advance() // consume TAG
		names, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicyUnsetTag
		stmt.UnsetTags = names

	case p.cur.Type == kwSET && p.peekNext().Type == kwBODY && d.returns != "":
		p.advance() // consume SET
		p.advance() // consume BODY
		if _, err := p.expect(tokArrow); err != nil {
			return nil, err
		}
		body, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicySetBody
		stmt.Body = body

	case p.cur.Type == kwSET:
		p.advance() // consume SET
		// FEATURE's SET may carry its AS '<definition>' alone; every other
		// SET names a property.
		need := !d.definition || d.kind == ast.PolicyIngress || p.cur.Type != kwAS
		opts, err := p.parsePolicyProps(d.set, need)
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicySet
		stmt.Options = opts
		if d.definition && p.cur.Type == kwAS {
			def, err := p.parsePolicyDefinition()
			if err != nil {
				return nil, err
			}
			stmt.Definition = def
		}

	case p.cur.Type == kwAS && d.definition:
		def, err := p.parsePolicyDefinition()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicyAs
		stmt.Definition = def

	case p.cur.Type == kwUNSET:
		p.advance() // consume UNSET
		unsetStart := p.cur
		keys, err := p.parseObjectUnsetKeys()
		if err != nil {
			return nil, err
		}
		if len(keys) > 1 && !d.unsetList {
			return nil, syntaxErrorAtTok(unsetStart)
		}
		for _, k := range keys {
			if !containsString(d.unset, k) {
				return nil, &ParseError{
					Position: unsetStart.Loc.Start, End: p.prev.Loc.End,
					Message: "syntax error: " + d.word + " POLICY has no property " + k + " to unset",
				}
			}
		}
		stmt.Action = ast.AlterPolicyUnset
		stmt.UnsetProps = keys

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropDocPolicy reads DROP <kind> POLICY [ IF EXISTS ] <name>, where
// d's DROP page writes IF EXISTS. cur is d's kind word; start is the DROP
// token's Loc.
func (p *Parser) parseDropDocPolicy(start ast.Loc, d *docPolicy) (ast.Node, error) {
	p.advance() // consume the kind word
	p.advance() // consume POLICY
	node, err := p.parseDropObject(ast.DropObject, start, d.dropIfExists, false)
	if err != nil {
		return nil, err
	}
	node.(*ast.DropStmt).Object = d.word + " POLICY"
	return node, nil
}

// parsePolicyProps reads a run of `KEY = value` properties, each one props
// documents, with a value of the shape it takes, and set once. need says
// whether an empty run is a syntax error.
func (p *Parser) parsePolicyProps(props policyProps, need bool) ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsObjectProperty() {
		opt, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if ok := props[opt.Name]; ok == nil || !ok(opt) {
			return nil, &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error at or near " + opt.Name}
		}
		if hasOption(opts, opt.Name) {
			return nil, &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error: " + opt.Name + " is set twice"}
		}
		opts = append(opts, opt)
	}
	if need && len(opts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return opts, nil
}

// parsePolicyDefinition reads AS '<definition>'. cur is AS.
func (p *Parser) parsePolicyDefinition() (*string, error) {
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	tok, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	s := tok.Str
	return &s, nil
}

// syntaxErrorAtTok reports a syntax error at tok, a token already consumed.
func syntaxErrorAtTok(tok Token) *ParseError {
	text := tok.Str
	if text == "" {
		text = TokenName(tok.Type)
	}
	return &ParseError{Position: tok.Loc.Start, End: tok.Loc.End, Message: "syntax error at or near " + text}
}

func hasOption(opts []*ast.CopyOption, name string) bool {
	for _, o := range opts {
		if o.Name == name {
			return true
		}
	}
	return false
}

func containsString(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}
