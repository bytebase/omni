// Package parser - endpoint.go implements T-SQL ENDPOINT statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// endpointStateValues: CREATE/ALTER ENDPOINT STATE = {STARTED | STOPPED | DISABLED}.
var endpointStateValues = newOptionSet().withIdents("STARTED", "STOPPED", "DISABLED")

// endpointAuthMethods: HTTP endpoint AUTHENTICATION options: NTLM/KERBEROS/NEGOTIATE.
var endpointAuthMethods = newOptionSet().withIdents("NTLM", "KERBEROS", "NEGOTIATE")

// endpointEncryptionAlgs: TCP endpoint ENCRYPTION ALGORITHM options: AES/RC4.
var endpointEncryptionAlgs = newOptionSet().withIdents("AES", "RC4")

// endpointRoleValues: service_broker role: WITNESS/PARTNER/ALL (kwALL exempted).
var endpointRoleValues = newOptionSet(kwALL).withIdents("WITNESS", "PARTNER", "ALL")

// endpointMessageForwarding: ENABLED/DISABLED.
var endpointMessageForwarding = newOptionSet().withIdents("ENABLED", "DISABLED")

// endpointPayloadTypes: first word of FOR clause. DATABASE and SERVICE are
// multi-word prefixes (DATABASE_MIRRORING, SERVICE_BROKER) handled via peek.
var endpointPayloadTypes = newOptionSet().withIdents(
	"TSQL", "SERVICE_BROKER", "DATABASE_MIRRORING", "SOAP",
	"SERVICE", "DATABASE",
)

// parseCreateEndpointStmt parses CREATE ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-endpoint-transact-sql
//
//	CREATE ENDPOINT endPointName [ AUTHORIZATION login ]
//	    [ STATE = { STARTED | STOPPED | DISABLED } ]
//	    AS { TCP } (
//	        <protocol_specific_arguments>
//	    )
//	    FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING } (
//	        <language_specific_arguments>
//	    )
//
//	<AS TCP_protocol_specific_arguments> ::=
//	AS TCP (
//	  LISTENER_PORT = listenerPort
//	  [ [ , ] LISTENER_IP = ALL | ( xx.xx.xx.xx IPv4 address ) | ( '__:__1' IPv6 address ) ]
//	)
//
//	<FOR TSQL_language_specific_arguments> ::=
//	FOR TSQL (
//	  [ ENCRYPTION = { NEGOTIATED | STRICT } ]
//	)
//
//	<FOR SERVICE_BROKER_language_specific_arguments> ::=
//	FOR SERVICE_BROKER (
//	   [ AUTHENTICATION = {
//	            WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	      | CERTIFICATE certificate_name
//	      | WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ] CERTIFICATE certificate_name
//	      | CERTIFICATE certificate_name WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	    } ]
//	   [ [ , ] ENCRYPTION = { DISABLED | { { SUPPORTED | REQUIRED }
//	       [ ALGORITHM { AES | RC4 | AES RC4 | RC4 AES } ] }
//	   ]
//	   [ [ , ] MESSAGE_FORWARDING = { ENABLED | DISABLED } ]
//	   [ [ , ] MESSAGE_FORWARD_SIZE = forward_size ]
//	)
//
//	<FOR DATABASE_MIRRORING_language_specific_arguments> ::=
//	FOR DATABASE_MIRRORING (
//	   [ AUTHENTICATION = {
//	            WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	      | CERTIFICATE certificate_name
//	      | WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ] CERTIFICATE certificate_name
//	      | CERTIFICATE certificate_name WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	   } ]
//	   [ [ , ] ENCRYPTION = { DISABLED | { { SUPPORTED | REQUIRED }
//	       [ ALGORITHM { AES | RC4 | AES RC4 | RC4 AES } ] }
//	   ]
//	   [ , ] ROLE = { WITNESS | PARTNER | ALL }
//	)
func (p *Parser) parseCreateEndpointStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// ENDPOINT keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "ENDPOINT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseEndpointOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterEndpointStmt parses ALTER ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-endpoint-transact-sql
//
//	ALTER ENDPOINT endPointName [ AUTHORIZATION login ]
//	    [ STATE = { STARTED | STOPPED | DISABLED } ]
//	    [ AS { TCP } (
//	        <protocol_specific_arguments>
//	    ) ]
//	    [ FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING } (
//	        <language_specific_arguments>
//	    ) ]
func (p *Parser) parseAlterEndpointStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// ENDPOINT keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "ENDPOINT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseEndpointOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropEndpointStmt parses DROP ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-endpoint-transact-sql
//
//	DROP ENDPOINT endPointName
func (p *Parser) parseDropEndpointStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// ENDPOINT keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "ENDPOINT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// endpointTopOptions defines the valid top-level endpoint options that are NOT
// handled by explicit switch cases (AUTHORIZATION, STATE, AS, FOR are handled
// explicitly). Only AFFINITY is a valid additional top-level option per the
// SqlScriptDOM grammar (endpointOption = endpointState | endpointAffinity).
var endpointTopOptions = newOptionSet(kwAFFINITY)

// parseEndpointOptions parses AUTHORIZATION, STATE, AS <protocol>, FOR <payload> clauses
// for CREATE/ALTER ENDPOINT statements.
//
//	[ AUTHORIZATION login ]
//	[ STATE = { STARTED | STOPPED | DISABLED } ]
//	[ AFFINITY = { NONE | ADMIN | integer } ]
//	[ AS { TCP } ( <protocol_specific_arguments> ) ]
//	[ FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING } ( <language_specific_arguments> ) ]
func (p *Parser) parseEndpointOptions() *nodes.List {
	var opts []nodes.Node

	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		switch {
		case p.cur.Type == kwAUTHORIZATION:
			optLoc := p.pos()
			p.advance()
			if p.isIdentLike() {
				opts = append(opts, &nodes.EndpointOption{Name: "AUTHORIZATION", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}

		case p.cur.Type == kwSTATE:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// STATE = { STARTED | STOPPED | DISABLED }
			if p.isValidOption(endpointStateValues) {
				opts = append(opts, &nodes.EndpointOption{Name: "STATE", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}

		case p.cur.Type == kwAS:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == kwTCP {
				opts = append(opts, &nodes.EndpointOption{Name: "AS", Value: "TCP", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
				p.parseEndpointProtocolTCPOptions(&opts)
			} else if p.cur.Type == kwHTTP {
				opts = append(opts, &nodes.EndpointOption{Name: "AS", Value: "HTTP", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
				p.parseEndpointProtocolHTTPOptions(&opts)
			} else if p.isAnyKeywordIdent() {
				// Unknown protocol - parse as generic key=value options
				opts = append(opts, &nodes.EndpointOption{Name: "AS", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
				p.parseEndpointGenericProtocolOptions(&opts)
			}

		case p.cur.Type == kwFOR:
			optLoc := p.pos()
			p.advance()
			// FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING | SOAP }
			// SqlScriptDOM's EndpointPayloadType enum; SERVICE / DATABASE are
			// multi-word prefixes handled by the lookahead below.
			if p.isValidOption(endpointPayloadTypes) {
				payloadType := strings.ToUpper(p.cur.Str)
				// Check for multi-word payload types
				if payloadType == "SERVICE" {
					p.advance()
					if p.cur.Type == kwBROKER {
						payloadType = "SERVICE_BROKER"
						p.advance()
					}
				} else if payloadType == "DATABASE" {
					p.advance()
					if p.cur.Type == kwMIRRORING {
						payloadType = "DATABASE_MIRRORING"
						p.advance()
					}
				} else {
					// TSQL or other single-word type
					p.advance()
				}
				opts = append(opts, &nodes.EndpointOption{Name: "FOR", Value: payloadType, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.parseEndpointPayloadOptions(payloadType, &opts)
			}

		default:
			// Check for additional valid top-level endpoint options (e.g., AFFINITY).
			// Reject any keyword/identifier not in the valid set to prevent accepting
			// invalid options like SELECT = STARTED.
			if p.isValidOption(endpointTopOptions) {
				optLoc := p.pos()
				key := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
						val := strings.ToUpper(p.cur.Str)
						if p.cur.Type == tokSCONST {
							val = "'" + p.cur.Str + "'"
						} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
							val = p.cur.Str
						}
						opts = append(opts, &nodes.EndpointOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
						p.advance()
					}
				} else {
					opts = append(opts, &nodes.EndpointOption{Name: key, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				}
			} else if p.isAnyKeywordIdent() {
				// Unknown keyword — not a valid endpoint option; stop parsing.
				goto done
			} else if p.cur.Type == ',' {
				// Comma between endpoint options (STATE = X, AFFINITY = Y)
				p.advance()
			} else {
				p.advance()
			}
		}
	}

done:
	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// parseEndpointProtocolTCPOptions parses TCP protocol-specific arguments.
//
//	( LISTENER_PORT = listenerPort
//	  [ [ , ] LISTENER_IP = ALL | ( xx.xx.xx.xx ) | ( 'ipv6' ) ] )
func (p *Parser) parseEndpointProtocolTCPOptions(opts *[]nodes.Node) {
	if p.cur.Type != '(' {
		return
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		if p.cur.Type == kwLISTENER_PORT {
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
				*opts = append(*opts, &nodes.EndpointOption{Name: "LISTENER_PORT", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else if p.cur.Type == kwLISTENER_IP {
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwALL {
				*opts = append(*opts, &nodes.EndpointOption{Name: "LISTENER_IP", Value: "ALL", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			} else if p.cur.Type == '(' {
				// ( ip_address ) or ( 'ipv6' )
				p.advance()
				if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
					ip := p.parseEndpointIPAddress()
					*opts = append(*opts, &nodes.EndpointOption{Name: "LISTENER_IP", Value: ip, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				}
				p.match(')')
			} else if p.isIdentLike() || p.cur.Type == tokSCONST {
				*opts = append(*opts, &nodes.EndpointOption{Name: "LISTENER_IP", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		} else {
			p.advance()
		}
	}

	p.match(')') // consume ')'
}

// parseEndpointGenericProtocolOptions parses protocol options for unknown/future protocols
// as generic key=value pairs.
//
//	( key = value [ , key = value ] ... )
func (p *Parser) parseEndpointGenericProtocolOptions(opts *[]nodes.Node) {
	if p.cur.Type != '(' {
		return
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		if p.isAnyKeywordIdent() {
			optLoc := p.pos()
			key := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
				if p.cur.Type == '(' {
					// Parenthesized value list: ( val1, val2, ... )
					p.advance() // consume '('
					var vals []string
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
							val := strings.ToUpper(p.cur.Str)
							if p.cur.Type == tokSCONST {
								val = "'" + p.cur.Str + "'"
							} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
								val = p.cur.Str
							}
							vals = append(vals, val)
							p.advance()
						} else {
							p.advance()
						}
					}
					p.match(')') // consume ')'
					*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: strings.Join(vals, ","), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				} else if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
					val := strings.ToUpper(p.cur.Str)
					if p.cur.Type == tokSCONST {
						val = "'" + p.cur.Str + "'"
					} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
						val = p.cur.Str
					}
					*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
					p.advance()
				}
			} else {
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			}
		} else {
			p.advance()
		}
	}

	p.match(')') // consume ')'
}

// parseEndpointProtocolHTTPOptions parses HTTP protocol-specific arguments (deprecated since SQL Server 2012).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-endpoint-transact-sql (pre-2012)
//
//	( PATH = 'url'
//	  , AUTHENTICATION = ( { BASIC | DIGEST | NTLM | KERBEROS | INTEGRATED } [ , ...n ] )
//	  , PORTS = ( { CLEAR | SSL } [ , ...n ] )
//	  [ , SITE = { '*' | '+' | 'webSite' } ]
//	  [ , CLEAR_PORT = clearPort ]
//	  [ , SSL_PORT = sslPort ]
//	  [ , AUTH_REALM = { 'realm' | NONE } ]
//	  [ , DEFAULT_LOGON_DOMAIN = { 'domain' | NONE } ]
//	  [ , COMPRESSION = { ENABLED | DISABLED } ]
//	)
func (p *Parser) parseEndpointProtocolHTTPOptions(opts *[]nodes.Node) {
	if p.cur.Type != '(' {
		return
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		if !p.isAnyKeywordIdent() {
			p.advance()
			continue
		}

		optLoc := p.pos()
		key := strings.ToUpper(p.cur.Str)
		p.advance()

		if p.cur.Type != '=' {
			*opts = append(*opts, &nodes.EndpointOption{Name: key, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			continue
		}
		p.advance() // consume '='

		switch key {
		case "AUTHENTICATION", "PORTS":
			// These take a parenthesized list: ( VALUE [, VALUE ...] )
			if p.cur.Type == '(' {
				p.advance()
				var vals []string
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
						vals = append(vals, strings.ToUpper(p.cur.Str))
						p.advance()
					} else {
						p.advance()
					}
				}
				p.match(')') // consume ')'
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: strings.Join(vals, ","), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
			} else if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		default:
			// Simple KEY = VALUE (string, number, or identifier)
			if p.cur.Type == tokSCONST {
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: "'" + p.cur.Str + "'", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			} else if p.cur.Type == tokICONST {
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			} else if p.isAnyKeywordIdent() {
				*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}
		}
	}

	p.match(')') // consume ')'
}

// parseEndpointIPAddress parses an IP address which may contain dots (parsed as separate tokens).
func (p *Parser) parseEndpointIPAddress() string {
	var parts []string
	if p.cur.Type == tokSCONST {
		// Quoted IPv6 like '::1'
		s := p.cur.Str
		p.advance()
		return "'" + s + "'"
	}
	// IPv4 like 10.0.75.1 — tokens: 10 . 0 . 75 . 1
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
			parts = append(parts, p.cur.Str)
			p.advance()
		} else if p.cur.Type == '.' {
			parts = append(parts, ".")
			p.advance()
		} else {
			break
		}
	}
	return strings.Join(parts, "")
}

// parseEndpointPayloadOptions parses payload-specific arguments within parentheses.
func (p *Parser) parseEndpointPayloadOptions(payloadType string, opts *[]nodes.Node) {
	if p.cur.Type != '(' {
		return
	}
	p.advance() // consume '('

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		switch {
		case p.cur.Type == kwAUTHENTICATION:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			authStr := p.parseEndpointAuthentication()
			*opts = append(*opts, &nodes.EndpointOption{Name: "AUTHENTICATION", Value: authStr, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

		case p.cur.Type == kwENCRYPTION:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			encStr := p.parseEndpointEncryption(payloadType)
			*opts = append(*opts, &nodes.EndpointOption{Name: "ENCRYPTION", Value: encStr, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

		case p.cur.Type == kwMESSAGE_FORWARDING:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// MESSAGE_FORWARDING = { ENABLED | DISABLED }
			if p.isValidOption(endpointMessageForwarding) {
				*opts = append(*opts, &nodes.EndpointOption{Name: "MESSAGE_FORWARDING", Value: strings.ToUpper(p.cur.Str), Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}

		case p.cur.Type == kwMESSAGE_FORWARD_SIZE:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == tokICONST || p.isAnyKeywordIdent() {
				*opts = append(*opts, &nodes.EndpointOption{Name: "MESSAGE_FORWARD_SIZE", Value: p.cur.Str, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}

		case p.cur.Type == kwROLE:
			optLoc := p.pos()
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			// ROLE = { WITNESS | PARTNER | ALL }
			if p.isValidOption(endpointRoleValues) {
				val := strings.ToUpper(p.cur.Str)
				*opts = append(*opts, &nodes.EndpointOption{Name: "ROLE", Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				p.advance()
			}

		default:
			// Unknown payload option as KEY = VALUE or bare keyword
			if p.isAnyKeywordIdent() {
				optLoc := p.pos()
				key := strings.ToUpper(p.cur.Str)
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
						val := strings.ToUpper(p.cur.Str)
						if p.cur.Type == tokSCONST {
							val = "'" + p.cur.Str + "'"
						} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
							val = p.cur.Str
						}
						*opts = append(*opts, &nodes.EndpointOption{Name: key, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
						p.advance()
					}
				} else {
					*opts = append(*opts, &nodes.EndpointOption{Name: key, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
				}
			} else {
				p.advance()
			}
		}
	}

	p.match(')') // consume ')'
}

// parseEndpointAuthentication parses the AUTHENTICATION value for SERVICE_BROKER and DATABASE_MIRRORING.
//
//	AUTHENTICATION = {
//	    WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	  | CERTIFICATE certificate_name
//	  | WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ] CERTIFICATE certificate_name
//	  | CERTIFICATE certificate_name WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	}
func (p *Parser) parseEndpointAuthentication() string {
	var parts []string

	if p.cur.Type == kwWINDOWS {
		parts = append(parts, "WINDOWS")
		p.advance()
		// Optional auth method: NTLM | KERBEROS | NEGOTIATE
		if p.isValidOption(endpointAuthMethods) {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		}
		// Optional CERTIFICATE after WINDOWS
		if p.cur.Type == kwCERTIFICATE {
			p.advance()
			if p.isIdentLike() {
				parts = append(parts, "CERTIFICATE", p.cur.Str)
				p.advance()
			}
		}
	} else if p.cur.Type == kwCERTIFICATE {
		p.advance()
		if p.isIdentLike() {
			parts = append(parts, "CERTIFICATE", p.cur.Str)
			p.advance()
		}
		// Optional WINDOWS after CERTIFICATE
		if p.cur.Type == kwWINDOWS {
			parts = append(parts, "WINDOWS")
			p.advance()
			if p.isValidOption(endpointAuthMethods) {
				parts = append(parts, strings.ToUpper(p.cur.Str))
				p.advance()
			}
		}
	}

	return strings.Join(parts, " ")
}

// parseEndpointEncryption parses the ENCRYPTION value.
//
// For TSQL payload:
//
//	ENCRYPTION = { NEGOTIATED | STRICT }
//
// For SERVICE_BROKER / DATABASE_MIRRORING:
//
//	ENCRYPTION = { DISABLED | { { SUPPORTED | REQUIRED }
//	    [ ALGORITHM { AES | RC4 | AES RC4 | RC4 AES } ] } }
func (p *Parser) parseEndpointEncryption(payloadType string) string {
	if !p.isAnyKeywordIdent() {
		return ""
	}

	val := strings.ToUpper(p.cur.Str)
	p.advance()

	// For TSQL: NEGOTIATED or STRICT - no further options
	if payloadType == "TSQL" {
		return val
	}

	// For SERVICE_BROKER / DATABASE_MIRRORING:
	// DISABLED has no further options
	if val == "DISABLED" {
		return val
	}

	// SUPPORTED or REQUIRED may have ALGORITHM clause
	if val == "SUPPORTED" || val == "REQUIRED" {
		if p.cur.Type == kwALGORITHM {
			p.advance()
			if p.isValidOption(endpointEncryptionAlgs) {
				alg1 := strings.ToUpper(p.cur.Str)
				p.advance()
				// Check for two-algorithm spec: AES RC4 or RC4 AES
				if p.isValidOption(endpointEncryptionAlgs) {
					alg2 := strings.ToUpper(p.cur.Str)
					if alg2 == "AES" || alg2 == "RC4" {
						p.advance()
						return val + " ALGORITHM " + alg1 + " " + alg2
					}
				}
				return val + " ALGORITHM " + alg1
			}
		}
	}

	return val
}
