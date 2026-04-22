// Package parser - security_misc.go implements T-SQL SECURITY POLICY,
// SENSITIVITY CLASSIFICATION, and SIGNATURE parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSecurityPolicyStmt parses a CREATE SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-security-policy-transact-sql
//
//	CREATE SECURITY POLICY [schema_name.] security_policy_name
//	    { ADD [ FILTER | BLOCK ] PREDICATE tvf_schema_name.security_predicate_function_name
//	      ( { column_name | expression } [ , ...n ] ) ON table_schema_name.table_name
//	      [ <block_dml_operation> ] } [ , ...n ]
//	    [ WITH ( STATE = { ON | OFF } [,] [ SCHEMABINDING = { ON | OFF } ] ) ]
//	    [ NOT FOR REPLICATION ]
//
//	<block_dml_operation>
//	    [ { AFTER { INSERT | UPDATE } }
//	    | { BEFORE { UPDATE | DELETE } } ]
func (p *Parser) parseCreateSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "CREATE",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// [schema_name.] policy_name
	stmt.Name , _ = p.parseTableRef()

	// Parse predicates
	stmt.Predicates = p.parseSecurityPredicateList()

	// WITH ( STATE = ON|OFF, SCHEMABINDING = ON|OFF )
	p.parseSecurityPolicyWithOptions(stmt)

	// NOT FOR REPLICATION
	if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwFOR {
			p.advance()
		}
		if p.cur.Type == kwREPLICATION {
			p.advance()
			stmt.NotForReplication = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterSecurityPolicyStmt parses an ALTER SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-security-policy-transact-sql
//
//	ALTER SECURITY POLICY schema_name.security_policy_name
//	    [
//	        { ADD { FILTER | BLOCK } PREDICATE tvf_schema_name.security_predicate_function_name
//	           ( { column_name | arguments } [ , ...n ] ) ON table_schema_name.table_name
//	           [ <block_dml_operation> ] }
//	        | { ALTER { FILTER | BLOCK } PREDICATE tvf_schema_name.new_security_predicate_function_name
//	             ( { column_name | arguments } [ , ...n ] ) ON table_schema_name.table_name
//	           [ <block_dml_operation> ] }
//	        | { DROP { FILTER | BLOCK } PREDICATE ON table_schema_name.table_name }
//	        | [ <additional_add_alter_drop_predicate_statements> [ , ...n ] ]
//	    ]    [ WITH ( STATE = { ON | OFF } ) ]
//	    [ NOT FOR REPLICATION ]
//
//	<block_dml_operation>
//	    [ { AFTER { INSERT | UPDATE } }
//	    | { BEFORE { UPDATE | DELETE } } ]
func (p *Parser) parseAlterSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "ALTER",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// schema_name.policy_name
	stmt.Name , _ = p.parseTableRef()

	// Parse predicates
	stmt.Predicates = p.parseSecurityPredicateList()

	// WITH ( STATE = ON|OFF )
	p.parseSecurityPolicyWithOptions(stmt)

	// NOT FOR REPLICATION
	if p.cur.Type == kwNOT {
		p.advance()
		if p.cur.Type == kwFOR {
			p.advance()
		}
		if p.cur.Type == kwREPLICATION {
			p.advance()
			stmt.NotForReplication = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropSecurityPolicyStmt parses a DROP SECURITY POLICY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-security-policy-transact-sql
//
//	DROP SECURITY POLICY [ IF EXISTS ] [schema_name.] security_policy_name
func (p *Parser) parseDropSecurityPolicyStmt() (*nodes.SecurityPolicyStmt, error) {
	loc := p.pos()
	// POLICY keyword already consumed by caller

	stmt := &nodes.SecurityPolicyStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		p.match(kwEXISTS)
		stmt.IfExists = true
	}

	// [schema_name.] policy_name
	stmt.Name , _ = p.parseTableRef()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSecurityPredicateList parses a comma-separated list of security predicates.
func (p *Parser) parseSecurityPredicateList() *nodes.List {
	var preds []nodes.Node
	for {
		pred := p.parseSecurityPredicate()
		if pred == nil {
			break
		}
		preds = append(preds, pred)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(preds) == 0 {
		return nil
	}
	return &nodes.List{Items: preds}
}

// parseSecurityPredicate parses a single ADD/ALTER/DROP FILTER/BLOCK PREDICATE clause.
func (p *Parser) parseSecurityPredicate() *nodes.SecurityPredicate {
	// Expect ADD, ALTER, or DROP
	var action string
	switch p.cur.Type {
	case kwADD:
		action = "ADD"
	case kwALTER:
		action = "ALTER"
	case kwDROP:
		action = "DROP"
	default:
		return nil
	}
	loc := p.pos()
	p.advance()

	pred := &nodes.SecurityPredicate{
		Action: action,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// FILTER or BLOCK
	if p.cur.Type == kwFILTER || p.cur.Type == kwBLOCK {
		pred.PredicateType = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// PREDICATE
	if p.cur.Type == kwPREDICATE {
		p.advance()
	}

	// For DROP: no function, just ON table
	if action == "DROP" {
		// ON table
		if p.cur.Type == kwON {
			p.advance()
			pred.Table , _ = p.parseTableRef()
		}
		pred.Loc.End = p.prevEnd()
		return pred
	}

	// Function name: schema.function_name
	pred.Function , _ = p.parseTableRef()

	// ( args )
	if p.cur.Type == '(' {
		p.advance()
		var args []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg, _ := p.parseExpr()
			if arg != nil {
				args = append(args, arg)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if len(args) > 0 {
			pred.Args = &nodes.List{Items: args}
		}
	}

	// ON table
	if p.cur.Type == kwON {
		p.advance()
		pred.Table , _ = p.parseTableRef()
	}

	// Optional block_dml_operation: AFTER INSERT|UPDATE, BEFORE UPDATE|DELETE
	if p.cur.Type == kwAFTER || p.cur.Type == kwBEFORE {
		timing := strings.ToUpper(p.cur.Str)
		p.advance()
		if p.cur.Type == kwINSERT || p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE {
			pred.BlockDMLOp = timing + " " + strings.ToUpper(p.cur.Str)
			p.advance()
		}
	}

	pred.Loc.End = p.prevEnd()
	return pred
}

// parseSecurityPolicyWithOptions parses WITH (STATE = ON|OFF, SCHEMABINDING = ON|OFF).
func (p *Parser) parseSecurityPolicyWithOptions(stmt *nodes.SecurityPolicyStmt) {
	if p.cur.Type != kwWITH {
		return
	}
	p.advance()
	if p.cur.Type != '(' {
		return
	}
	p.advance()

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isAnyKeywordIdent() {
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			switch key {
			case "STATE":
				if p.cur.Type == kwON {
					v := true
					stmt.StateOn = &v
					p.advance()
				} else if p.cur.Type == kwOFF {
					v := false
					stmt.StateOn = &v
					p.advance()
				}
			case "SCHEMABINDING":
				if p.cur.Type == kwON {
					v := true
					stmt.SchemaBinding = &v
					p.advance()
				} else if p.cur.Type == kwOFF {
					v := false
					stmt.SchemaBinding = &v
					p.advance()
				}
			default:
				if p.isAnyKeywordIdent() || p.cur.Type == kwON || p.cur.Type == kwOFF {
					p.advance()
				}
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	p.match(')')
}

// parseAddSensitivityClassificationStmt parses an ADD SENSITIVITY CLASSIFICATION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/add-sensitivity-classification-transact-sql
//
//	ADD SENSITIVITY CLASSIFICATION TO
//	    <object_name> [ , ...n ]
//	    WITH ( <sensitivity_option> [ , ...n ] )
//
//	<object_name> ::= [ schema_name. ] table_name.column_name
//
//	<sensitivity_option> ::=
//	    LABEL = string |
//	    LABEL_ID = guidOrString |
//	    INFORMATION_TYPE = string |
//	    INFORMATION_TYPE_ID = guidOrString |
//	    RANK = NONE | LOW | MEDIUM | HIGH | CRITICAL
func (p *Parser) parseAddSensitivityClassificationStmt() (*nodes.SensitivityClassificationStmt, error) {
	loc := p.pos()
	// CLASSIFICATION keyword already consumed by caller

	stmt := &nodes.SensitivityClassificationStmt{
		Action: "ADD",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// object_name list (comma-separated qualified names: schema.table.column)
	var cols []nodes.Node
	for {
		ref , _ := p.parseTableRef()
		if ref != nil {
			cols = append(cols, ref)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}

	// WITH ( options )
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			var opts []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isAnyKeywordIdent() {
					optLoc := p.pos()
					key := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					val := ""
					if p.cur.Type == tokSCONST {
						val = p.cur.Str
						p.advance()
					} else if p.isIdentLike() || p.cur.Type == tokICONST {
						val = p.cur.Str
						p.advance()
					}
					opts = append(opts, &nodes.SensitivityOption{
						Key:   key,
						Value: val,
						Loc:   nodes.Loc{Start: optLoc, End: p.prevEnd()},
					})
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.match(')')
			if len(opts) > 0 {
				stmt.Options = &nodes.List{Items: opts}
			}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropSensitivityClassificationStmt parses a DROP SENSITIVITY CLASSIFICATION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-sensitivity-classification-transact-sql
//
//	DROP SENSITIVITY CLASSIFICATION FROM
//	    <object_name> [ , ...n ]
//
//	<object_name> ::= [ schema_name. ] table_name.column_name
func (p *Parser) parseDropSensitivityClassificationStmt() (*nodes.SensitivityClassificationStmt, error) {
	loc := p.pos()
	// CLASSIFICATION keyword already consumed by caller

	stmt := &nodes.SensitivityClassificationStmt{
		Action: "DROP",
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// object_name list
	var cols []nodes.Node
	for {
		ref , _ := p.parseTableRef()
		if ref != nil {
			cols = append(cols, ref)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSignatureStmt parses ADD/DROP [COUNTER] SIGNATURE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/add-signature-transact-sql
//
//	ADD [ COUNTER ] SIGNATURE TO module_class::module_name
//	    BY <crypto_list> [ , ...n ]
//
//	DROP [ COUNTER ] SIGNATURE FROM module_class::module_name
//	    BY <crypto_list> [ , ...n ]
//
//	<crypto_list> ::=
//	    CERTIFICATE cert_name
//	    | CERTIFICATE cert_name [ WITH PASSWORD = 'password' ]
//	    | CERTIFICATE cert_name WITH SIGNATURE = signed_blob
//	    | ASYMMETRIC KEY Asym_Key_Name
//	    | ASYMMETRIC KEY Asym_Key_Name [ WITH PASSWORD = 'password' ]
//	    | ASYMMETRIC KEY Asym_Key_Name WITH SIGNATURE = signed_blob
func (p *Parser) parseSignatureStmt(action string) (*nodes.SignatureStmt, error) {
	loc := p.pos()
	// SIGNATURE keyword already consumed by caller

	stmt := &nodes.SignatureStmt{
		Action: action,
		Loc:    nodes.Loc{Start: loc, End: -1},
	}

	// TO (for ADD) or FROM (for DROP)
	if p.cur.Type == kwTO {
		p.advance()
	} else if p.cur.Type == kwFROM {
		p.advance()
	}

	// module_class::module_name  or just module_name
	// module_class is OBJECT (default), ASSEMBLY, DATABASE, etc.
	if p.isAnyKeywordIdent() {
		name1 := p.cur.Str
		next := p.peekNext()
		if next.Type == tokCOLONCOLON {
			// Single-word class:: pattern
			stmt.ModuleClass = strings.ToUpper(name1)
			p.advance() // consume class word
			p.advance() // consume ::
			stmt.ModuleName , _ = p.parseTableRef()
		} else {
			// Could be a dotted name or multi-word class (e.g., ASYMMETRIC KEY::)
			p.advance() // consume name1
			// Check if next is another word followed by ::
			if p.isAnyKeywordIdent() && p.peekNext().Type == tokCOLONCOLON {
				stmt.ModuleClass = strings.ToUpper(name1) + " " + strings.ToUpper(p.cur.Str)
				p.advance() // consume second word
				p.advance() // consume ::
				stmt.ModuleName , _ = p.parseTableRef()
			} else {
				// No :: — reconstruct as dotted name
				ref := &nodes.TableRef{Object: name1}
				for p.cur.Type == '.' {
					p.advance()
					if p.isIdentLike() {
						ref = &nodes.TableRef{
							Schema: ref.Object,
							Object: p.cur.Str,
						}
						p.advance()
					}
				}
				stmt.ModuleName = ref
			}
		}
	}

	// BY <crypto_list>
	if p.cur.Type == kwBY {
		p.advance()
		var cryptos []nodes.Node
		for {
			crypto := p.parseSignatureCryptoItem()
			if crypto == nil {
				break
			}
			cryptos = append(cryptos, crypto)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if len(cryptos) > 0 {
			stmt.CryptoList = &nodes.List{Items: cryptos}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSignatureCryptoItem parses a single crypto reference in a SIGNATURE BY clause.
//
//	<crypto_list> ::=
//	    CERTIFICATE cert_name
//	    | CERTIFICATE cert_name [ WITH PASSWORD = 'password' ]
//	    | CERTIFICATE cert_name WITH SIGNATURE = signed_blob
//	    | ASYMMETRIC KEY Asym_Key_Name
//	    | ASYMMETRIC KEY Asym_Key_Name [ WITH PASSWORD = 'password' ]
//	    | ASYMMETRIC KEY Asym_Key_Name WITH SIGNATURE = signed_blob
func (p *Parser) parseSignatureCryptoItem() *nodes.CryptoItem {
	loc := p.pos()
	item := &nodes.CryptoItem{Loc: nodes.Loc{Start: loc, End: -1}}

	// CERTIFICATE cert_name or ASYMMETRIC KEY key_name
	if p.cur.Type == kwCERTIFICATE {
		item.Mechanism = "CERTIFICATE"
		p.advance()
		if p.isIdentLike() {
			item.Name = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwASYMMETRIC {
		item.Mechanism = "ASYMMETRIC KEY"
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
		}
		if p.isIdentLike() {
			item.Name = p.cur.Str
			p.advance()
		}
	} else {
		return nil
	}

	// Optional: WITH PASSWORD = 'password' or WITH SIGNATURE = hex_blob
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == kwPASSWORD {
			item.WithType = "PASSWORD"
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokSCONST {
				item.WithValue = p.cur.Str
				p.advance()
			}
		} else if p.cur.Type == kwSIGNATURE {
			item.WithType = "SIGNATURE"
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// hex blob
			if p.isIdentLike() || p.cur.Type == tokICONST {
				item.WithValue = p.cur.Str
				p.advance()
			}
		}
	}

	item.Loc.End = p.prevEnd()
	return item
}
