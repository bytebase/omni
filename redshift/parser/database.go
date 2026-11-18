package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// parseCreatedbStmt parses a CREATE DATABASE statement.
// The CREATE keyword has already been consumed. Current token is DATABASE.
//
// Ref: gram.y CreatedbStmt
//
//	CREATE DATABASE name opt_with createdb_opt_list
func (p *Parser) parseCreatedbStmt(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Redshift's data-source clause, which PostgreSQL has no form of:
	// FROM INTEGRATION, FROM ARN, and [WITH PERMISSIONS] FROM DATASHARE.
	if p.isRedshiftCreateDatabaseSourceStart() {
		return p.parseRedshiftCreateDatabaseSource(name, stmtLoc)
	}

	// opt_with: WITH | EMPTY
	if p.cur.Type == WITH {
		p.advance()
	}

	options, err := p.parseCreatedbOptList()
	if err != nil {
		return nil, err
	}

	return &nodes.CreatedbStmt{
		Dbname:  name,
		Options: options,
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseRedshiftCreateDatabaseSource parses CREATE DATABASE's Redshift data
// source: the FROM clause of r_CREATE_DATABASE, after the database name and,
// for a datashare, an optional WITH PERMISSIONS.
//
//	FROM INTEGRATION '<id>' [ DATABASE '<name>' ] [ SET ] <integration options>
//	    [ WITH ] <createdb options>
//	| FROM ARN '<arn>' ( WITH DATA CATALOG SCHEMA '<schema>'
//	    | WITH NO DATA CATALOG SCHEMA )
//	| [ WITH PERMISSIONS ] FROM DATASHARE <share> OF [ ACCOUNT <id> ]
//	    NAMESPACE <id>
func (p *Parser) parseRedshiftCreateDatabaseSource(dbname string, stmtLoc int) (nodes.Node, error) {
	var items []nodes.Node
	withPermissions := false
	if p.cur.Type == WITH {
		p.advance() // WITH PERMISSIONS
		if !p.acceptRedshiftKeyword("permissions") {
			return nil, p.syntaxErrorAtCur()
		}
		withPermissions = true
		items = append(items, makeDefElem("with_permissions", &nodes.Boolean{Boolval: true}))
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}

	switch {
	case withPermissions && !p.redshiftKeywordIs("datashare"):
		return nil, p.syntaxErrorAtCur()
	case p.acceptRedshiftKeyword("integration"):
		id, err := p.parseRedshiftSourceString()
		if err != nil {
			return nil, err
		}
		items = append(items, makeDefElem("from_integration", &nodes.String{Str: id}))
		if p.acceptRedshiftKeyword("database") {
			name, err := p.parseRedshiftSourceName()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("integration_database", &nodes.String{Str: name}))
		}
		if p.acceptRedshiftKeyword("set") {
			items = append(items, makeDefElem("set", &nodes.Boolean{Boolval: true}))
		}
		options, err := p.parseRedshiftIntegrationOptions(false, redshiftCreateDatabaseIntegrationOrder)
		if err != nil {
			return nil, err
		}
		items = append(items, options...)
		if p.cur.Type == WITH {
			p.advance()
		}
		createdbOptions, err := p.parseRedshiftIntegrationCreatedbOptions()
		if err != nil {
			return nil, err
		}
		if createdbOptions != nil {
			items = append(items, createdbOptions.Items...)
		}
	case p.acceptRedshiftKeyword("arn"):
		arn, err := p.parseRedshiftSourceString()
		if err != nil {
			return nil, err
		}
		items = append(items, makeDefElem("from_arn", &nodes.String{Str: arn}))
		if _, err := p.expect(WITH); err != nil {
			return nil, err
		}
		noData := false
		if p.acceptRedshiftKeyword("no") {
			noData = true
		}
		for _, word := range []string{"data", "catalog", "schema"} {
			if !p.acceptRedshiftKeyword(word) {
				return nil, p.syntaxErrorAtCur()
			}
		}
		if noData {
			items = append(items, makeDefElem("no_data_catalog_schema", &nodes.Boolean{Boolval: true}))
		} else {
			schema, err := p.parseRedshiftSourceString()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("data_catalog_schema", &nodes.String{Str: schema}))
		}
		// The page's examples carry an IAM role after the catalog clause.
		if p.acceptRedshiftKeyword("iam_role") {
			role := ""
			switch {
			case p.acceptRedshiftKeyword("default"):
				role = "default"
			case p.cur.Type == SCONST:
				role = p.advance().Str
			default:
				return nil, p.syntaxErrorAtCur()
			}
			items = append(items, makeDefElem("iam_role", &nodes.String{Str: role}))
		}
	case p.acceptRedshiftKeyword("datashare"):
		share, err := p.parseName()
		if err != nil {
			return nil, err
		}
		items = append(items, makeDefElem("from_datashare", &nodes.String{Str: share}))
		if _, err := p.expect(OF); err != nil {
			return nil, err
		}
		if p.acceptRedshiftKeyword("account") {
			account, err := p.parseRedshiftDatashareValue()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("datashare_account", &nodes.String{Str: account}))
		}
		if !p.acceptRedshiftKeyword("namespace") {
			return nil, p.syntaxErrorAtCur()
		}
		namespace, err := p.parseRedshiftDatashareValue()
		if err != nil {
			return nil, err
		}
		items = append(items, makeDefElem("datashare_namespace", &nodes.String{Str: namespace}))
	default:
		return nil, p.syntaxErrorAtCur()
	}

	return &nodes.CreatedbStmt{
		Dbname:  dbname,
		Options: &nodes.List{Items: items},
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseRedshiftSourceString reads the string literal a Redshift data-source
// clause requires.
func (p *Parser) parseRedshiftSourceString() (string, error) {
	if p.cur.Type != SCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// parseRedshiftSourceName reads an integration source database's name: the
// string literal the Syntax block writes or the identifier its examples use.
func (p *Parser) parseRedshiftSourceName() (string, error) {
	if p.cur.Type == SCONST {
		return p.advance().Str, nil
	}
	return p.parseName()
}

// parseRedshiftQualifiedNameStringList parses a comma-separated list of
// qualified names into one dotted String per name.
func (p *Parser) parseRedshiftQualifiedNameStringList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		parts, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(parts.Items))
		for _, part := range parts.Items {
			names = append(names, part.(*nodes.String).Str)
		}
		items = append(items, &nodes.String{Str: strings.Join(names, ".")})
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseRedshiftDatashareValue reads a datashare account or namespace: a name,
// a number, or the string literal the AWS pages write those values as.
func (p *Parser) parseRedshiftDatashareValue() (string, error) {
	switch p.cur.Type {
	case SCONST, ICONST:
		return p.advance().Str, nil
	default:
		return p.parseName()
	}
}

// parseRedshiftIntegrationCreatedbOptions parses what r_CREATE_DATABASE
// allows after an integration source: the owner, connection limit, collation
// and isolation level, in that order, each at most once. The generic
// createdb_opt_list would also take any identifier as an option name.
func (p *Parser) parseRedshiftIntegrationCreatedbOptions() (*nodes.List, error) {
	var items []nodes.Node
	next := 0
	for {
		idx := -1
		switch p.cur.Type {
		case OWNER:
			idx = 0
		case CONNECTION:
			idx = 1
		case COLLATE:
			idx = 2
		case ISOLATION:
			idx = 3
		}
		if idx < 0 {
			break
		}
		if idx < next {
			return nil, p.syntaxErrorAtCur()
		}
		next = idx + 1
		item, err := p.parseCreatedbOptItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// redshiftIntegrationOptionIndex is the position of the current option in
// order, or -1.
func (p *Parser) redshiftIntegrationOptionIndex(order []string) int {
	for i, word := range order {
		if p.redshiftKeywordIs(word) {
			return i
		}
	}
	return -1
}

// The two pages order the integration options differently: r_CREATE_DATABASE
// starts at ACCEPTINVCHARS, r_ALTER_DATABASE's SET QUERY_ALL_STATES form starts
// at QUERY_ALL_STATES and carries ACCEPTINVCHARS after it.
var (
	redshiftCreateDatabaseIntegrationOrder = []string{
		"acceptinvchars", "query_all_states", "refresh_interval",
		"truncatecolumns", "history_mode", "auto_remediation",
	}
	redshiftAlterDatabaseIntegrationOrder = []string{
		"query_all_states", "acceptinvchars", "refresh_interval",
		"truncatecolumns", "history_mode", "auto_remediation",
	}
)

// parseRedshiftIntegrationOptions parses the option list shared by
// CREATE DATABASE ... FROM INTEGRATION and ALTER DATABASE ... INTEGRATION SET.
// historyModeFor is true in ALTER, whose HISTORY_MODE takes a FOR clause; order
// is that statement's documented sequence, each option at most once.
func (p *Parser) parseRedshiftIntegrationOptions(historyModeFor bool, order []string) ([]nodes.Node, error) {
	var items []nodes.Node
	next := 0
	for {
		idx := p.redshiftIntegrationOptionIndex(order)
		if idx < 0 {
			break
		}
		if idx < next {
			return nil, p.syntaxErrorAtCur()
		}
		next = idx + 1
		switch {
		case p.acceptRedshiftKeyword("acceptinvchars"):
			value, err := p.parseRedshiftIntegrationBool()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("acceptinvchars", &nodes.Boolean{Boolval: value}))
		case p.acceptRedshiftKeyword("query_all_states"):
			value, err := p.parseRedshiftIntegrationBool()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("query_all_states", &nodes.Boolean{Boolval: value}))
		case p.acceptRedshiftKeyword("refresh_interval"):
			if p.cur.Type != ICONST {
				return nil, p.syntaxErrorAtCur()
			}
			items = append(items, makeDefElem("refresh_interval", &nodes.Integer{Ival: p.parseSignedIconst()}))
		case p.acceptRedshiftKeyword("truncatecolumns"):
			value, err := p.parseRedshiftIntegrationBool()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("truncatecolumns", &nodes.Boolean{Boolval: value}))
		case p.acceptRedshiftKeyword("history_mode"):
			value, err := p.parseRedshiftIntegrationBool()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("history_mode", &nodes.Boolean{Boolval: value}))
			if !historyModeFor || !p.acceptRedshiftKeyword("for") {
				continue
			}
			switch {
			case p.acceptRedshiftKeyword("all"):
				if _, err := p.expect(TABLES); err != nil {
					return nil, err
				}
				items = append(items, makeDefElem("history_mode_all_tables", &nodes.Boolean{Boolval: true}))
				if p.acceptRedshiftKeyword("in") {
					if _, err := p.expect(SCHEMA); err != nil {
						return nil, err
					}
					schemas, err := p.parseNameList()
					if err != nil {
						return nil, err
					}
					items = append(items, makeDefElem("history_mode_schemas", schemas))
				}
			case p.acceptRedshiftKeyword("table"):
				tables, err := p.parseRedshiftQualifiedNameStringList()
				if err != nil {
					return nil, err
				}
				items = append(items, makeDefElem("history_mode_tables", tables))
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case p.acceptRedshiftKeyword("auto_remediation"):
			value, err := p.parseRedshiftIntegrationBool()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("auto_remediation", &nodes.Boolean{Boolval: value}))
		}
	}
	return items, nil
}

// parseRedshiftIntegrationBool parses [=] TRUE|FALSE.
func (p *Parser) parseRedshiftIntegrationBool() (bool, error) {
	if p.cur.Type == '=' {
		p.advance()
	}
	switch p.cur.Type {
	case TRUE_P:
		p.advance()
		return true, nil
	case FALSE_P:
		p.advance()
		return false, nil
	default:
		return false, p.syntaxErrorAtCur()
	}
}

// parseCreatedbOptList parses createdb_opt_list.
//
//	createdb_opt_list:
//	    createdb_opt_items | EMPTY
func (p *Parser) parseCreatedbOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseCreatedbOptItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseCreatedbOptItem parses a single createdb_opt_item.
//
//	createdb_opt_item:
//	    createdb_opt_name opt_equal NumericOnly
//	    | createdb_opt_name opt_equal opt_boolean_or_string
//	    | createdb_opt_name opt_equal DEFAULT
func (p *Parser) parseCreatedbOptItem() (*nodes.DefElem, error) {
	optName, err := p.parseCreatedbOptName()
	if err != nil || optName == "" {
		return nil, err
	}

	// opt_equal: '=' | EMPTY
	if p.cur.Type == '=' {
		p.advance()
	}

	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		return makeDefElem(optName, nil), nil
	case ICONST, FCONST, '+', '-':
		val, err := p.parseNumericOnlyRequired()
		if err != nil {
			return nil, err
		}
		return makeDefElem(optName, val), nil
	case TRUE_P, FALSE_P, ON, SCONST:
		return makeDefElem(optName, &nodes.String{Str: p.parseOptBooleanOrString()}), nil
	}
	// opt_boolean_or_string: ... | NonReservedWord_or_Sconst
	if p.isNonReservedWord() {
		return makeDefElem(optName, &nodes.String{Str: p.parseOptBooleanOrString()}), nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseCreatedbOptName parses createdb_opt_name.
//
//	createdb_opt_name:
//	    IDENT
//	    | CONNECTION LIMIT
//	    | ENCODING
//	    | LOCATION
//	    | OWNER
//	    | TABLESPACE
//	    | TEMPLATE
func (p *Parser) parseCreatedbOptName() (string, error) {
	switch p.cur.Type {
	case COLLATE:
		p.advance()
		return "collate", nil
	case CONNECTION:
		p.advance() // consume CONNECTION
		if _, err := p.expect(LIMIT); err != nil {
			return "", err
		}
		return "connection_limit", nil
	case ENCODING:
		p.advance()
		return "encoding", nil
	case LOCATION:
		p.advance()
		return "location", nil
	case OWNER:
		p.advance()
		return "owner", nil
	case ISOLATION:
		p.advance()
		p.expect(LEVEL)
		return "isolation_level", nil
	case TABLESPACE:
		p.advance()
		return "tablespace", nil
	case TEMPLATE:
		p.advance()
		return "template", nil
	case IDENT:
		name := p.cur.Str
		p.advance()
		return name, nil
	default:
		return "", nil
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

// parseAlterDatabaseDispatch parses ALTER DATABASE statements.
// The ALTER keyword has already been consumed. Current token is DATABASE.
//
// This dispatches between:
//   - AlterDatabaseStmt (ALTER DATABASE name [WITH] createdb_opt_list)
//   - AlterDatabaseStmt (ALTER DATABASE name SET TABLESPACE name)
//   - AlterDatabaseSetStmt (ALTER DATABASE name SET/RESET ...)
//   - AlterDatabaseStmt (ALTER DATABASE name CONNECTION LIMIT ...)
func (p *Parser) parseAlterDatabaseDispatch(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	if p.isRedshiftAlterDatabaseIntegrationStart() {
		return p.parseRedshiftAlterDatabaseIntegration(name, stmtLoc)
	}

	switch p.cur.Type {
	case RENAME:
		// ALTER DATABASE name RENAME TO name
		p.advance() // consume RENAME
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_DATABASE,
			Subname:    name,
			Newname:    newname,
		}, nil

	case OWNER:
		// ALTER DATABASE name OWNER TO RoleSpec; without TO, OWNER is a
		// createdb_opt_name (OWNER [=] value).
		if p.peekNext().Type != TO && !p.completingAtNext() {
			options, err := p.parseCreatedbOptList()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterDatabaseStmt{
				Dbname:  name,
				Options: options,
				Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
			}, nil
		}
		p.advance() // consume OWNER
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec, err := p.parseRoleSpec()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_DATABASE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil

	case SET:
		// Could be SET TABLESPACE or SET variable
		next := p.peekNext()
		if next.Type == TABLESPACE {
			// ALTER DATABASE name SET TABLESPACE name
			p.advance() // consume SET
			p.advance() // consume TABLESPACE
			tbsName, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterDatabaseStmt{
				Dbname:  name,
				Options: &nodes.List{Items: []nodes.Node{makeDefElem("tablespace", &nodes.String{Str: tbsName})}},
				Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
			}, nil
		}
		// ALTER DATABASE name SET variable = value (SetResetClause)
		return p.parseAlterDatabaseSetStmt(name, stmtLoc)

	case RESET:
		// ALTER DATABASE name RESET ... (SetResetClause)
		return p.parseAlterDatabaseSetStmt(name, stmtLoc)

	case WITH:
		// ALTER DATABASE name WITH createdb_opt_list
		p.advance() // consume WITH
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	case CONNECTION:
		// ALTER DATABASE name CONNECTION LIMIT connlimit
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	case REFRESH:
		// ALTER DATABASE name REFRESH COLLATION VERSION
		p.advance() // consume REFRESH
		if _, err := p.expect(COLLATION); err != nil {
			return nil, err
		}
		if _, err := p.expect(VERSION_P); err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseRefreshCollStmt{
			Dbname: name,
			Loc:    nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	default:
		// ALTER DATABASE name createdb_opt_list (empty or options)
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil
	}
}

// parseRedshiftAlterDatabaseIntegration parses Redshift's ALTER DATABASE
// INTEGRATION clause from r_ALTER_DATABASE, after the database name.
//
//	INTEGRATION REFRESH
//	    ( { ALL | INERROR | REMEDIABLE } TABLES [ IN SCHEMA <name> [, ...] ]
//	    | TABLE <name> [, ...] )
//	| INTEGRATION SET <integration options>
func (p *Parser) parseRedshiftAlterDatabaseIntegration(dbname string, stmtLoc int) (nodes.Node, error) {
	p.advance() // INTEGRATION
	var items []nodes.Node
	switch {
	case p.acceptRedshiftKeyword("refresh"):
		mode := ""
		switch {
		case p.acceptRedshiftKeyword("all"):
			mode = "all"
		case p.acceptRedshiftKeyword("inerror"):
			mode = "inerror"
		case p.acceptRedshiftKeyword("remediable"):
			mode = "remediable"
		case p.acceptRedshiftKeyword("table"):
			mode = "table"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, makeDefElem("refresh", &nodes.String{Str: mode}))
		if mode == "table" {
			tables, err := p.parseRedshiftQualifiedNameStringList()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("refresh_tables", tables))
			break
		}
		if _, err := p.expect(TABLES); err != nil {
			return nil, err
		}
		if p.acceptRedshiftKeyword("in") {
			if _, err := p.expect(SCHEMA); err != nil {
				return nil, err
			}
			schemas, err := p.parseNameList()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("refresh_schemas", schemas))
		}
	case p.acceptRedshiftKeyword("set"):
		options, err := p.parseRedshiftIntegrationOptions(true, redshiftAlterDatabaseIntegrationOrder)
		if err != nil {
			return nil, err
		}
		if len(options) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, options...)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.AlterDatabaseStmt{
		Dbname:  dbname,
		Options: &nodes.List{Items: items},
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseAlterDatabaseSetStmt parses ALTER DATABASE name SET/RESET ... .
// The name has already been parsed. Current token is SET or RESET.
func (p *Parser) parseAlterDatabaseSetStmt(dbname string, stmtLoc int) (nodes.Node, error) {
	setstmt, err := p.parseSetResetClause()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterDatabaseSetStmt{
		Dbname:  dbname,
		Setstmt: setstmt,
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseSetResetClause parses SetResetClause: SET set_rest, which takes
// no LOCAL or SESSION, or RESET reset_rest.
func (p *Parser) parseSetResetClause() (*nodes.VariableSetStmt, error) {
	var result nodes.Node
	var err error
	switch p.cur.Type {
	case SET:
		setLoc := p.pos()
		p.advance()
		if result, err = p.parseSetRest(); err != nil {
			return nil, err
		}
		if vs, ok := result.(*nodes.VariableSetStmt); ok {
			vs.Loc = nodes.Loc{Start: setLoc, End: p.prev.End}
		}
	case RESET:
		p.advance()
		if result, err = p.parseVariableResetStmt(); err != nil {
			return nil, err
		}
	default:
		if p.collectMode() {
			p.addTokenCandidate(SET)
			p.addTokenCandidate(RESET)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	vs, _ := result.(*nodes.VariableSetStmt)
	return vs, nil
}

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

// parseDropdbStmt parses a DROP DATABASE statement.
// The DROP keyword has already been consumed. Current token is DATABASE.
//
// Ref: gram.y DropdbStmt
//
//	DROP DATABASE name
//	| DROP DATABASE IF_P EXISTS name
//	| DROP DATABASE name opt_with '(' drop_option_list ')'
//	| DROP DATABASE IF_P EXISTS name opt_with '(' drop_option_list ')'
func (p *Parser) parseDropdbStmt(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance() // consume IF
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Redshift's DROP DATABASE database_name [ FORCE ] is a bare keyword, and it
	// stands where PostgreSQL's opt_with '(' drop_option_list ')' does: the two
	// are alternatives, not a prefix and a suffix.
	var options *nodes.List
	switch {
	case p.cur.Type == FORCE:
		p.advance()
		options = appendToList(options, makeDefElem("force", nil))
	case p.cur.Type == WITH || p.cur.Type == '(':
		if p.cur.Type == WITH {
			p.advance() // consume WITH
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if options, err = p.parseDropOptionList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	return &nodes.DropdbStmt{
		Dbname:    name,
		MissingOk: missingOk,
		Options:   options,
		Loc:       nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseDropOptionList parses drop_option_list.
//
//	drop_option_list:
//	    drop_option
//	    | drop_option_list ',' drop_option
func (p *Parser) parseDropOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt, err := p.parseDropOption()
		if err != nil {
			return nil, err
		}
		items = append(items, opt)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseDropOption parses drop_option.
//
//	drop_option:
//	    FORCE
func (p *Parser) parseDropOption() (*nodes.DefElem, error) {
	if _, err := p.expect(FORCE); err != nil {
		return nil, err
	}
	return makeDefElem("force", nil), nil
}
