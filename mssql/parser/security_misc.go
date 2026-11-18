// Package parser - security_misc.go implements T-SQL SECURITY POLICY,
// SENSITIVITY CLASSIFICATION, and SIGNATURE parsing.
package parser

import (
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSecurityPolicyStmt parses a CREATE SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-security-policy-transact-sql
// Grammar: TSql160.g:7639 createSecurityPolicyStatement_SECURITY
//
//	CREATE SECURITY POLICY schemaObjectThreePartName
//	    ( createSecurityPredicateListElement ( Comma createSecurityPredicateListElement )* )?
//	    ( With LeftParenthesis securityPolicyOptionsList RightParenthesis )?
//	    ( Not For Replication )?
func (p *Parser) parseCreateSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "CREATE",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// schema_name.policy_name: the grammar names one and SQL Server refuses
	// `CREATE SECURITY POLICY` alone (Msg 102, measured).
	name, err := p.parseSecurityPolicyName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Predicates: only ADD leads a CREATE's list (createSecurityPredicateListElement).
	stmt.Predicates, err = p.parseSecurityPredicateList("CREATE")
	if err != nil {
		return nil, err
	}

	if err := p.parseSecurityPolicyWithOptions(stmt, "CREATE"); err != nil {
		return nil, err
	}

	// NOT FOR REPLICATION
	if p.cur.Type == kwNOT {
		p.advance()
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwREPLICATION); err != nil {
			return nil, err
		}
		stmt.NotForReplication = true
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterSecurityPolicyStmt parses an ALTER SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-security-policy-transact-sql
// Grammar: TSql160.g:8711 alterSecurityPolicyStatement_SECURITY
//
//	ALTER SECURITY POLICY schemaObjectThreePartName
//	    (   ( With LeftParenthesis securityPolicyOptionsList RightParenthesis )
//	      | ( alterSecurityPredicateListElement ( Comma alterSecurityPredicateListElement )* )
//	      | Add Not For Replication
//	      | Drop Not For Replication )
//
// The body is one alternative, not a sequence: SQL Server refuses a predicate
// list followed by WITH (Msg 319) and ADD NOT FOR REPLICATION followed by WITH
// (Msg 319), both measured.
func (p *Parser) parseAlterSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "ALTER",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// schema_name.policy_name
	name, err := p.parseSecurityPolicyName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.cur.Type == kwWITH:
		err = p.parseSecurityPolicyWithOptions(stmt, "ALTER")
	case p.cur.Type == kwADD && p.peekNext().Type == kwNOT,
		p.cur.Type == kwDROP && p.peekNext().Type == kwNOT:
		// ADD NOT FOR REPLICATION sets the flag and DROP NOT FOR REPLICATION
		// clears it (alterSecurityPolicyStatement, TSql160.g:23439), so the two
		// statements read differently.
		add := p.cur.Type == kwADD
		p.advance() // ADD or DROP
		p.advance() // NOT
		if _, e := p.expect(kwFOR); e != nil {
			err = e
			break
		}
		if _, e := p.expect(kwREPLICATION); e != nil {
			err = e
			break
		}
		stmt.NotForReplication = add
	case p.cur.Type == kwADD, p.cur.Type == kwALTER, p.cur.Type == kwDROP:
		stmt.Predicates, err = p.parseSecurityPredicateList("ALTER")
	default:
		err = p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropSecurityPolicyStmt parses a DROP SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-security-policy-transact-sql
// Grammar: TSql160.g:9015 dropSecurityPolicyStatement_SECURITY
//
//	DROP SECURITY POLICY ( If Exists )? dropObjectList
func (p *Parser) parseDropSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// IF EXISTS: IF takes EXISTS, and SQL Server refuses `DROP SECURITY POLICY
	// IF a` and a bare `IF` (Msg 102, measured).
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	// [schema_name.] policy_name
	name, err := p.parseSecurityPolicyName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityPolicyName reads the policy's schemaObjectThreePartName and
// refuses a statement that names none.
func (p *Parser) parseSecurityPolicyName() (*nodes.TableRef, error) {
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return ref, nil
}

// parseSecurityPredicateList parses a comma-separated list of security
// predicates. action is the statement's: CREATE takes only ADD, ALTER takes
// ADD, ALTER and DROP (alterSecurityPredicateListElement, TSql160.g:4907).
func (p *Parser) parseSecurityPredicateList(action string) (*nodes.List, error) {
	var preds []nodes.Node
	for {
		pred, err := p.parseSecurityPredicate(action)
		if err != nil {
			return nil, err
		}
		if pred == nil {
			break
		}
		preds = append(preds, pred)
		if _, ok := p.match(','); !ok {
			break
		}
		// A comma joins two predicates: `ADD FILTER ... ON a, DROP ...` is no
		// statement (SQL Server reads the DROP as another statement, Msg 343).
		if !p.startsSecurityPredicate(action) {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if len(preds) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: preds}, nil
}

// startsSecurityPredicate reports whether the current token can begin a
// predicate of a statement whose action is action (CREATE or ALTER).
func (p *Parser) startsSecurityPredicate(action string) bool {
	switch p.cur.Type {
	case kwADD:
		return true
	case kwALTER, kwDROP:
		return action == "ALTER"
	}
	return false
}

// parseSecurityPredicate parses one ADD/ALTER/DROP FILTER|BLOCK PREDICATE
// clause, or answers nil when the current token begins none.
//
// Grammar: TSql160.g:4925 createOrAltersecurityPredicateListElement and
// TSql160.g:4930 dropSecurityPredicateListElement. The predicate type is an
// Identifier the actions match against FILTER and BLOCK, so `ADD PREDICATE`
// is no predicate (Msg 102, measured).
func (p *Parser) parseSecurityPredicate(action string) (*nodes.SecurityPredicate, error) {
	loc := p.pos()
	var predAction string
	switch p.cur.Type {
	case kwADD:
		predAction = "ADD"
	case kwALTER:
		if action != "ALTER" {
			return nil, nil
		}
		predAction = "ALTER"
	case kwDROP:
		if action != "ALTER" {
			return nil, nil
		}
		predAction = "DROP"
	default:
		return nil, nil
	}
	p.advance()

	pred := &nodes.SecurityPredicate{
		Action: predAction,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// FILTER or BLOCK
	switch p.cur.Type {
	case kwFILTER:
		pred.PredicateType = "FILTER"
	case kwBLOCK:
		pred.PredicateType = "BLOCK"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if _, err := p.expect(kwPREDICATE); err != nil {
		return nil, err
	}

	// DROP carries no function: DROP { FILTER | BLOCK } PREDICATE ON table.
	if predAction != "DROP" {
		// schema.function_name
		fn, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if fn == nil {
			return nil, p.syntaxErrorAtCur()
		}
		pred.Function = fn

		// userFunctionCall: ( { column_name | arguments } [ , ...n ] )
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var args []nodes.Node
		for p.cur.Type != ')' {
			arg, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.syntaxErrorAtCur()
			}
			args = append(args, arg)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if len(args) > 0 {
			pred.Args = &nodes.List{Items: args}
		}
	}

	// ON table_schema_name.table_name
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if table == nil {
		return nil, p.syntaxErrorAtCur()
	}
	pred.Table = table

	if err := p.parseSecurityPredicateOperation(pred); err != nil {
		return nil, err
	}

	pred.Loc.End = p.prevEnd()
	return pred, nil
}

// parseSecurityPredicateOperation parses <block_dml_operation>, the optional
// AFTER { INSERT | UPDATE } or BEFORE { UPDATE | DELETE } tail.
//
// Grammar: TSql160.g:4938 securityPredicateOperation. Only a BLOCK predicate
// takes one: SQL Server refuses `ADD FILTER PREDICATE ... AFTER INSERT` with
// Msg 102 at the keyword (measured).
func (p *Parser) parseSecurityPredicateOperation(pred *nodes.SecurityPredicate) error {
	if p.cur.Type != kwAFTER && p.cur.Type != kwBEFORE {
		return nil
	}
	at := p.pos()
	timing := strings.ToUpper(p.cur.Str)
	p.advance()

	var op string
	switch p.cur.Type {
	case kwINSERT:
		op = "INSERT"
	case kwUPDATE:
		op = "UPDATE"
	case kwDELETE:
		op = "DELETE"
	default:
		return p.syntaxErrorAtCur()
	}
	if (timing == "AFTER" && op == "DELETE") || (timing == "BEFORE" && op == "INSERT") {
		return p.syntaxErrorAtCur()
	}
	p.advance()

	if pred.PredicateType == "FILTER" {
		return p.newParseError(at, "a FILTER predicate takes no BEFORE/AFTER operation")
	}
	pred.BlockDMLOp = timing + " " + op
	return nil
}

// parseSecurityPolicyWithOptions parses WITH (
// { STATE = ON|OFF [ , SCHEMABINDING = ON|OFF ] } ).
//
// Grammar: TSql160.g:4914 securityPolicyOptionsList. SQL Server holds
// SCHEMABINDING to CREATE (the option's action checks the statement's action):
// `ALTER SECURITY POLICY a WITH(SCHEMABINDING = ON)` is Msg 102, measured. A
// repeated option is Msg 1039, and ALTER takes one option only — two are
// Msg 102 at the comma.
func (p *Parser) parseSecurityPolicyWithOptions(stmt *nodes.SecurityPolicyStmt, action string) error {
	if p.cur.Type != kwWITH {
		return nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return err
	}

	seen := map[string]bool{}
	for {
		at := p.pos()
		var key string
		switch p.cur.Type {
		case kwSTATE:
			key = "STATE"
		case kwSCHEMABINDING:
			key = "SCHEMABINDING"
			if action != "CREATE" {
				return p.syntaxErrorAtCur()
			}
		default:
			return p.syntaxErrorAtCur()
		}
		if seen[key] {
			return p.newParseError(at, "Option '"+key+"' is specified more than once.")
		}
		seen[key] = true
		p.advance()

		if _, err := p.expect('='); err != nil {
			return err
		}
		var value bool
		switch p.cur.Type {
		case kwON:
			value = true
		case kwOFF:
			value = false
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if key == "STATE" {
			v := value
			stmt.StateOn = &v
		} else {
			v := value
			stmt.SchemaBinding = &v
		}

		if p.cur.Type == ',' && action != "CREATE" {
			// ALTER takes one option: a second is Msg 102 at the comma.
			return p.syntaxErrorAtCur()
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}

	_, err := p.expect(')')
	return err
}

// sensitivityOptionTypes are the option names SensitivityClassification
// knows: any other is "Option 'x' is not supported for sensitivity
// classification." (Msg 16101, measured).
var sensitivityOptionTypes = []string{
	"LABEL",
	"LABEL_ID",
	"INFORMATION_TYPE",
	"INFORMATION_TYPE_ID",
	"RANK",
}

// sensitivityRanks are the words RANK takes: any other is Msg 16114 (measured).
var sensitivityRanks = []string{
	"NONE",
	"LOW",
	"MEDIUM",
	"HIGH",
	"CRITICAL",
}

// parseAddSensitivityClassificationStmt parses an ADD SENSITIVITY CLASSIFICATION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/add-sensitivity-classification-transact-sql
// Grammar: TSql160.g:8819 addSensitivityClassificationStatement_SENSITIVITY
//
//	ADD SENSITIVITY CLASSIFICATION TO column ( Comma column )*
//	    With LeftParenthesis sensitivityClassificationOption
//	    ( Comma sensitivityClassificationOption )* RightParenthesis
func (p *Parser) parseAddSensitivityClassificationStmt() (*nodes.SensitivityClassificationStmt, error) {
	loc := p.pos()
	// CLASSIFICATION keyword already consumed by caller

	stmt := &nodes.SensitivityClassificationStmt{
		Action: "ADD",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// TO
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}

	columns, err := p.parseSensitivityColumns()
	if err != nil {
		return nil, err
	}
	stmt.Columns = columns

	// WITH ( options ), which ADD does not make optional.
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	options, err := p.parseSensitivityOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = options
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropSensitivityClassificationStmt parses a DROP SENSITIVITY CLASSIFICATION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-sensitivity-classification-transact-sql
// Grammar: TSql160.g:9021 dropSensitivityClassificationStatement_SENSITIVITY
//
//	DROP SENSITIVITY CLASSIFICATION From column ( Comma column )*
func (p *Parser) parseDropSensitivityClassificationStmt() (*nodes.SensitivityClassificationStmt, error) {
	loc := p.pos()
	// CLASSIFICATION keyword already consumed by caller

	stmt := &nodes.SensitivityClassificationStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// FROM
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	columns, err := p.parseSensitivityColumns()
	if err != nil {
		return nil, err
	}
	stmt.Columns = columns

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSensitivityColumns parses the comma-separated columns ADD and DROP
// SENSITIVITY CLASSIFICATION name.
func (p *Parser) parseSensitivityColumns() (*nodes.List, error) {
	var cols []nodes.Node
	for {
		ref, err := p.parseSensitivityColumn()
		if err != nil {
			return nil, err
		}
		cols = append(cols, ref)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: cols}, nil
}

// parseSensitivityColumn parses one `<object_name>`: [schema_name.]
// table_name.column_name (the statement's docs).
//
// SQL Server refuses a column with no table name (Msg 16103), one with a
// database part (Msg 16110) and one with an empty part (Msg 1004), all
// measured, so the name holds two or three parts and none empty.
func (p *Parser) parseSensitivityColumn() (*nodes.TableRef, error) {
	loc := p.pos()
	parts, err := p.parseMultiPartIdentifier("table_ref")
	if err != nil {
		return nil, err
	}
	if len(parts) < 2 || len(parts) > 3 {
		return nil, p.newParseError(loc, "a sensitivity classification names a column as [schema_name.]table_name.column_name")
	}
	for _, part := range parts {
		if part == "" {
			return nil, p.newParseError(loc, "a sensitivity classification names a column as [schema_name.]table_name.column_name")
		}
	}
	ref := &nodes.TableRef{Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
	assignObjectRefParts(ref, parts)
	return ref, nil
}

// parseSensitivityOptions parses the option list of ADD SENSITIVITY
// CLASSIFICATION.
//
// Grammar: TSql160.g:4317 sensitivityClassificationOption, with the checks its
// action makes: the name is one of the five option types (Msg 16101), a string
// is the value of every option but RANK (Msg 16115) and RANK takes one of the
// five ranks (Msgs 16113, 16114), no option repeats (Msg 16104) and no string
// is empty (Msg 16106). All measured on the pinned server.
func (p *Parser) parseSensitivityOptions() (*nodes.List, error) {
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		optLoc := p.pos()
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		key := strings.ToUpper(name)
		if !slices.Contains(sensitivityOptionTypes, key) {
			return nil, p.newParseError(optLoc, "Option '"+key+"' is not supported for sensitivity classification.")
		}
		if seen[key] {
			return nil, p.newParseError(optLoc, "Sensitivity classification option '"+key+"' was repeated.")
		}
		seen[key] = true

		if _, err := p.expect('='); err != nil {
			return nil, err
		}

		value := ""
		if key == "RANK" {
			rank, ok := p.parseIdentifier()
			if !ok {
				return nil, p.newParseError(optLoc, "Sensitivity rank must not be provided as string.")
			}
			if !slices.Contains(sensitivityRanks, strings.ToUpper(rank)) {
				return nil, p.newParseError(optLoc, "Sensitivity rank must be one of the following: NONE, LOW, MEDIUM, HIGH, CRITICAL.")
			}
			value = rank
		} else {
			if p.cur.Type != tokSCONST {
				if p.isIdentLike() {
					return nil, p.newParseError(p.pos(), "'"+key+"' option's value must be provided as string.")
				}
				return nil, p.syntaxErrorAtCur()
			}
			value = p.cur.Str
			if value == "" {
				return nil, p.newParseError(p.pos(), "Option '"+key+"' is empty.")
			}
			p.advance()
		}

		opts = append(opts, &nodes.SensitivityOption{
			Key:   key,
			Value: value,
			Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
		})

		if _, ok := p.match(','); !ok {
			break
		}
	}
	return &nodes.List{Items: opts}, nil
}
