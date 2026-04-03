// Package parser - availability.go implements T-SQL AVAILABILITY GROUP statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateAvailabilityGroupStmt parses CREATE AVAILABILITY GROUP.
// Caller has consumed CREATE AVAILABILITY GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-availability-group-transact-sql
//
//	CREATE AVAILABILITY GROUP group_name
//	   WITH (<with_option_spec> [ ,...n ] )
//	   FOR [ DATABASE database_name [ ,...n ] ]
//	   REPLICA ON <add_replica_spec> [ ,...n ]
//	   AVAILABILITY GROUP ON <add_availability_group_spec> [ ,...2 ]
//	   [ LISTENER 'dns_name' ( <listener_option> ) ]
//	[ ; ]
//
//	<with_option_spec>::=
//	    AUTOMATED_BACKUP_PREFERENCE = { PRIMARY | SECONDARY_ONLY| SECONDARY | NONE }
//	  | FAILURE_CONDITION_LEVEL  = { 1 | 2 | 3 | 4 | 5 }
//	  | HEALTH_CHECK_TIMEOUT = milliseconds
//	  | DB_FAILOVER  = { ON | OFF }
//	  | DTC_SUPPORT  = { PER_DB | NONE }
//	  | [ BASIC | DISTRIBUTED | CONTAINED [ REUSE_SYSTEM_DATABASES | AUTOSEEDING_SYSTEM_DATABASES ] ]
//	  | REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT = { integer }
//	  | CLUSTER_TYPE = { WSFC | EXTERNAL | NONE }
//	  | WRITE_LEASE_VALIDITY = { seconds }
//	  | CLUSTER_CONNECTION_OPTIONS = 'key_value_pairs'
//
//	<add_replica_spec>::=
//	  <server_instance> WITH
//	    (
//	       ENDPOINT_URL = 'TCP://system-address:port',
//	       AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT | CONFIGURATION_ONLY },
//	       FAILOVER_MODE = { AUTOMATIC | MANUAL | EXTERNAL }
//	       [ , <add_replica_option> [ ,...n ] ]
//	    )
//
//	  <add_replica_option>::=
//	       SEEDING_MODE = { AUTOMATIC | MANUAL }
//	     | BACKUP_PRIORITY = n
//	     | SECONDARY_ROLE ( { [ ALLOW_CONNECTIONS = { NO | READ_ONLY | ALL } ]
//	        [,] [ READ_ONLY_ROUTING_URL = 'TCP://system-address:port' ] } )
//	     | PRIMARY_ROLE ( { [ ALLOW_CONNECTIONS = { READ_WRITE | ALL } ]
//	        [,] [ READ_ONLY_ROUTING_LIST = { ( '<server_instance>' [ ,...n ] ) | NONE } ]
//	        [,] [ READ_WRITE_ROUTING_URL = 'TCP://system-address:port' ] } )
//	     | SESSION_TIMEOUT = integer
//
//	<add_availability_group_spec>::=
//	 <ag_name> WITH
//	    (
//	       LISTENER_URL = 'TCP://system-address:port',
//	       AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT },
//	       FAILOVER_MODE = MANUAL,
//	       SEEDING_MODE = { AUTOMATIC | MANUAL }
//	    )
//
//	<listener_option> ::=
//	   {
//	      WITH DHCP [ ON ( <network_subnet_option> ) ]
//	    | WITH IP ( { ( <ip_address_option> ) } [ , ...n ] ) [ , PORT = listener_port ]
//	   }
//
//	  <network_subnet_option> ::= 'ip4_address', 'four_part_ipv4_mask'
//	  <ip_address_option> ::= { 'ip4_address', 'pv4_mask' | 'ipv6_address' }
func (p *Parser) parseCreateAvailabilityGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// AVAILABILITY GROUP already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name
	if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAvailabilityGroupOptions()

	// Reject if tokens remain unconsumed (e.g., invalid option names like SELECT).
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterAvailabilityGroupStmt parses ALTER AVAILABILITY GROUP.
// Caller has consumed ALTER AVAILABILITY GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-availability-group-transact-sql
//
//	ALTER AVAILABILITY GROUP group_name
//	  {
//	     SET ( <set_option_spec> )
//	   | ADD DATABASE database_name
//	   | REMOVE DATABASE database_name
//	   | ADD REPLICA ON <add_replica_spec>
//	   | MODIFY REPLICA ON <modify_replica_spec>
//	   | REMOVE REPLICA ON <server_instance>
//	   | JOIN
//	   | JOIN AVAILABILITY GROUP ON <add_availability_group_spec> [ , ...2 ]
//	   | MODIFY AVAILABILITY GROUP ON <modify_availability_group_spec> [ , ...2 ]
//	   | GRANT CREATE ANY DATABASE
//	   | DENY CREATE ANY DATABASE
//	   | FAILOVER
//	   | FORCE_FAILOVER_ALLOW_DATA_LOSS
//	   | ADD LISTENER 'dns_name' ( <add_listener_option> )
//	   | MODIFY LISTENER 'dns_name' ( <modify_listener_option> )
//	   | RESTART LISTENER 'dns_name'
//	   | REMOVE LISTENER 'dns_name'
//	   | OFFLINE
//	  }
//	[ ; ]
//
//	<set_option_spec> ::=
//	    AUTOMATED_BACKUP_PREFERENCE = { PRIMARY | SECONDARY_ONLY | SECONDARY | NONE }
//	  | FAILURE_CONDITION_LEVEL  = { 1 | 2 | 3 | 4 | 5 }
//	  | HEALTH_CHECK_TIMEOUT = milliseconds
//	  | DB_FAILOVER  = { ON | OFF }
//	  | DTC_SUPPORT  = { PER_DB | NONE }
//	  | REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT = { integer }
//	  | ROLE = SECONDARY
//	  | CLUSTER_CONNECTION_OPTIONS = 'key_value_pairs'
//
//	<add_replica_spec>::=
//	  <server_instance> WITH
//	    (
//	       ENDPOINT_URL = 'TCP://system-address:port' ,
//	       AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT | CONFIGURATION_ONLY } ,
//	       FAILOVER_MODE = { AUTOMATIC | MANUAL }
//	       [ , <add_replica_option> [ , ...n ] ]
//	    )
//
//	<modify_replica_spec>::=
//	  <server_instance> WITH
//	    (
//	       ENDPOINT_URL = 'TCP://system-address:port'
//	     | AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT }
//	     | FAILOVER_MODE = { AUTOMATIC | MANUAL }
//	     | SEEDING_MODE = { AUTOMATIC | MANUAL }
//	     | BACKUP_PRIORITY = n
//	     | SECONDARY_ROLE ( { ... } )
//	     | PRIMARY_ROLE ( { ... } )
//	     | SESSION_TIMEOUT = seconds
//	    )
//
//	<add_listener_option> ::=
//	   {
//	      WITH DHCP [ ON ( <network_subnet_option> ) ]
//	    | WITH IP ( { ( <ip_address_option> ) } [ , ...n ] ) [ , PORT = listener_port ]
//	   }
//
//	<modify_listener_option>::=
//	    {
//	       ADD IP ( <ip_address_option> )
//	     | PORT = listener_port
//	     | REMOVE IP ( 'ipv4_address' | 'ipv6_address')
//	    }
func (p *Parser) parseAlterAvailabilityGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// AVAILABILITY GROUP already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name
	if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Options = p.parseAvailabilityGroupOptions()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropAvailabilityGroupStmt parses DROP AVAILABILITY GROUP.
// Caller has consumed DROP AVAILABILITY GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-availability-group-transact-sql
//
//	DROP AVAILABILITY GROUP group_name
//	[ ; ]
func (p *Parser) parseDropAvailabilityGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// AVAILABILITY GROUP already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	// group_name
	if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// newAGOpt creates an AvailabilityGroupOption with the given name, value, and location.
func newAGOpt(name, value string, loc nodes.Loc) *nodes.AvailabilityGroupOption {
	return &nodes.AvailabilityGroupOption{Name: name, Value: value, Loc: loc}
}

// parseAvailabilityGroupOptions consumes the remaining clauses of
// CREATE/ALTER AVAILABILITY GROUP: WITH, FOR DATABASE, REPLICA ON,
// AVAILABILITY GROUP ON, LISTENER, SET, ADD DATABASE, REMOVE DATABASE,
// JOIN, FAILOVER, FORCE_FAILOVER_ALLOW_DATA_LOSS, OFFLINE, etc.
func (p *Parser) parseAvailabilityGroupOptions() *nodes.List {
	var opts []nodes.Node

	for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
		start := p.pos()
		switch {
		case p.cur.Type == '(':
			// Consume parenthesized blocks structurally (SET options, replica options, listener options, etc.)
			subOpts := p.parseParenthesizedAGOptions()
			opts = append(opts, subOpts...)

		case p.cur.Type == kwWITH:
			p.advance()
			opts = append(opts, newAGOpt("WITH", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwFOR:
			p.advance()
			if p.cur.Type == kwDATABASE {
				p.advance()
				var dbs []string
				for p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					dbs = append(dbs, p.cur.Str)
					p.advance()
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				opts = append(opts, newAGOpt("FOR DATABASE", strings.Join(dbs, ", "), nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("FOR", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwREPLICA:
			p.advance() // consume REPLICA
			if p.cur.Type == kwON {
				p.advance() // consume ON
			}
			opts = append(opts, newAGOpt("REPLICA ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwAVAILABILITY:
			p.advance() // consume AVAILABILITY
			if p.cur.Type == kwGROUP {
				p.advance() // consume GROUP
			}
			if p.cur.Type == kwON {
				p.advance() // consume ON
			}
			opts = append(opts, newAGOpt("AVAILABILITY GROUP ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwLISTENER:
			p.advance() // consume LISTENER
			name := p.consumeAGListenerName()
			opts = append(opts, newAGOpt("LISTENER", p.agQuoteName(name), nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwSET:
			p.advance()
			opts = append(opts, newAGOpt("SET", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwADD:
			p.advance() // consume ADD
			if p.cur.Type == kwDATABASE {
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					dbName := p.cur.Str
					p.advance()
					opts = append(opts, newAGOpt("ADD DATABASE", dbName, nodes.Loc{Start: start, End: p.prevEnd()}))
				}
			} else if p.cur.Type == kwREPLICA {
				p.advance() // consume REPLICA
				if p.cur.Type == kwON {
					p.advance() // consume ON
				}
				opts = append(opts, newAGOpt("ADD REPLICA ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else if p.cur.Type == kwLISTENER {
				p.advance() // consume LISTENER
				name := p.consumeAGListenerName()
				opts = append(opts, newAGOpt("ADD LISTENER", p.agQuoteName(name), nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("ADD", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwREMOVE:
			p.advance() // consume REMOVE
			if p.cur.Type == kwDATABASE {
				p.advance()
				if p.isAnyKeywordIdent() || p.cur.Type == tokSCONST {
					dbName := p.cur.Str
					p.advance()
					opts = append(opts, newAGOpt("REMOVE DATABASE", dbName, nodes.Loc{Start: start, End: p.prevEnd()}))
				}
			} else if p.cur.Type == kwREPLICA {
				p.advance() // consume REPLICA
				if p.cur.Type == kwON {
					p.advance() // consume ON
				}
				opts = append(opts, newAGOpt("REMOVE REPLICA ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else if p.cur.Type == kwLISTENER {
				p.advance() // consume LISTENER
				name := p.consumeAGListenerName()
				opts = append(opts, newAGOpt("REMOVE LISTENER", p.agQuoteName(name), nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("REMOVE", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwMODIFY:
			p.advance() // consume MODIFY
			if p.cur.Type == kwREPLICA {
				p.advance() // consume REPLICA
				if p.cur.Type == kwON {
					p.advance() // consume ON
				}
				opts = append(opts, newAGOpt("MODIFY REPLICA ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else if p.cur.Type == kwLISTENER {
				p.advance() // consume LISTENER
				name := p.consumeAGListenerName()
				opts = append(opts, newAGOpt("MODIFY LISTENER", p.agQuoteName(name), nodes.Loc{Start: start, End: p.prevEnd()}))
			} else if p.cur.Type == kwAVAILABILITY {
				p.advance() // consume AVAILABILITY
				if p.cur.Type == kwGROUP {
					p.advance() // consume GROUP
				}
				if p.cur.Type == kwON {
					p.advance() // consume ON
				}
				opts = append(opts, newAGOpt("MODIFY AVAILABILITY GROUP ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("MODIFY", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwJOIN:
			p.advance() // consume JOIN
			// JOIN AVAILABILITY GROUP ON ...
			if p.cur.Type == kwAVAILABILITY {
				p.advance() // consume AVAILABILITY
				if p.cur.Type == kwGROUP {
					p.advance() // consume GROUP
				}
				if p.cur.Type == kwON {
					p.advance() // consume ON
				}
				opts = append(opts, newAGOpt("JOIN AVAILABILITY GROUP ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("JOIN", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwGRANT:
			p.advance() // consume GRANT
			p.consumeAGCreateAnyDatabase()
			opts = append(opts, newAGOpt("GRANT CREATE ANY DATABASE", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwDENY:
			p.advance() // consume DENY
			p.consumeAGCreateAnyDatabase()
			opts = append(opts, newAGOpt("DENY CREATE ANY DATABASE", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwFAILOVER:
			p.advance()
			opts = append(opts, newAGOpt("FAILOVER", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwFORCE_FAILOVER_ALLOW_DATA_LOSS:
			p.advance()
			opts = append(opts, newAGOpt("FORCE_FAILOVER_ALLOW_DATA_LOSS", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwOFFLINE:
			p.advance()
			opts = append(opts, newAGOpt("OFFLINE", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == kwRESTART:
			p.advance() // consume RESTART
			if p.cur.Type == kwLISTENER {
				p.advance() // consume LISTENER
				name := p.consumeAGListenerName()
				opts = append(opts, newAGOpt("RESTART LISTENER", p.agQuoteName(name), nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt("RESTART", "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		case p.cur.Type == kwON:
			p.advance()
			opts = append(opts, newAGOpt("ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))

		case p.cur.Type == ',':
			p.advance() // skip commas

		case p.isAGOptionToken():
			key := strings.ToUpper(p.cur.Str)
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				key = "'" + p.cur.Str + "'"
			}
			p.advance()
			if p.cur.Type == '=' {
				p.advance() // consume '='
				if p.isAGOptionToken() {
					val := strings.ToUpper(p.cur.Str)
					if p.cur.Type == tokSCONST {
						val = "'" + p.cur.Str + "'"
					}
					p.advance()
					opts = append(opts, newAGOpt(key, val, nodes.Loc{Start: start, End: p.prevEnd()}))
				} else {
					opts = append(opts, newAGOpt(key, "", nodes.Loc{Start: start, End: p.prevEnd()}))
				}
			} else {
				opts = append(opts, newAGOpt(key, "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}

		default:
			// Stop parsing on unrecognized tokens — let caller detect leftover tokens.
			if len(opts) == 0 {
				return nil
			}
			return &nodes.List{Items: opts}
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}

// agQuoteName wraps a listener name in single quotes if non-empty, for display consistency.
func (p *Parser) agQuoteName(name string) string {
	if name != "" {
		return "'" + name + "'"
	}
	return ""
}

// consumeAGListenerName consumes a string constant listener name if present.
// Returns the listener name, or empty string if no string constant follows.
func (p *Parser) consumeAGListenerName() string {
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		name := p.cur.Str
		p.advance()
		return name
	}
	return ""
}

// consumeAGCreateAnyDatabase explicitly consumes the tokens CREATE ANY DATABASE.
// This replaces the previous generic loop that blindly consumed all identifier/CREATE tokens.
func (p *Parser) consumeAGCreateAnyDatabase() {
	// Expect: CREATE ANY DATABASE
	if p.cur.Type == kwCREATE {
		p.advance()
	}
	if p.cur.Type == kwANY {
		p.advance()
	}
	if p.cur.Type == kwDATABASE {
		p.advance()
	}
}

// parseParenthesizedAGOptions parses a parenthesized block of availability group options
// into individual AvailabilityGroupOption nodes instead of capturing raw text.
//
// Handles:
//   - Simple key = value options: ENDPOINT_URL = 'TCP://...' → Name="ENDPOINT_URL" Value="'TCP://...'"
//   - Nested options: SECONDARY_ROLE (ALLOW_CONNECTIONS = READ_ONLY) → Name="SECONDARY_ROLE" Value="(ALLOW_CONNECTIONS=READ_ONLY)"
//   - Value lists: READ_ONLY_ROUTING_LIST = ('server1', 'server2') → Name="READ_ONLY_ROUTING_LIST" Value="('server1', 'server2')"
//   - Standalone keywords: DISTRIBUTED → Name="DISTRIBUTED" Value=""
//   - Nested IP tuples: IP (('addr', 'mask')) → Name="IP" Value="(('addr', 'mask'))"
func (p *Parser) parseParenthesizedAGOptions() []nodes.Node {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // consume '('

	var opts []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		start := p.pos()

		// Handle nested parenthesized tuples (e.g., IP address tuples like ('addr', 'mask'))
		if p.cur.Type == '(' {
			tupleStr := p.parseAGParenTuple()
			if tupleStr != "" {
				opts = append(opts, newAGOpt("", "("+tupleStr+")", nodes.Loc{Start: start, End: p.prevEnd()}))
			}
			continue
		}

		// Must be an identifier, string, number, or keyword
		if !p.isAGOptionToken() {
			return opts
		}

		key := p.agOptionTokenString()
		p.advance()

		switch {
		case p.cur.Type == '=':
			p.advance() // consume '='
			if p.cur.Type == '(' {
				// Value is a parenthesized list: key = ( ... )
				subOpts := p.parseParenthesizedAGOptions()
				var parts []string
				for _, o := range subOpts {
					parts = append(parts, agOptDisplayStr(o))
				}
				opts = append(opts, newAGOpt(key, "("+strings.Join(parts, ", ")+")", nodes.Loc{Start: start, End: p.prevEnd()}))
			} else if p.isAGOptionToken() {
				val := p.agOptionTokenString()
				p.advance()
				opts = append(opts, newAGOpt(key, val, nodes.Loc{Start: start, End: p.prevEnd()}))
			} else {
				opts = append(opts, newAGOpt(key, "", nodes.Loc{Start: start, End: p.prevEnd()}))
			}
		case p.cur.Type == '(':
			// Nested parens without '=': SECONDARY_ROLE(...), PRIMARY_ROLE(...), IP(...), ON(...)
			// Collapse into Name to distinguish from key=(list) format
			subOpts := p.parseParenthesizedAGOptions()
			var parts []string
			for _, o := range subOpts {
				parts = append(parts, agOptDisplayStr(o))
			}
			opts = append(opts, newAGOpt(key+"("+strings.Join(parts, ", ")+")", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		default:
			opts = append(opts, newAGOpt(key, "", nodes.Loc{Start: start, End: p.prevEnd()}))
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return opts
}

// agOptDisplayStr returns the display string for an AG option node,
// reconstructing the original key=value format.
// Nested parens are collapsed into Name (e.g. "SECONDARY_ROLE(...)").
// Key=value pairs use Name + "=" + Value.
func agOptDisplayStr(n nodes.Node) string {
	switch opt := n.(type) {
	case *nodes.AvailabilityGroupOption:
		if opt.Value == "" {
			return opt.Name
		}
		if opt.Name == "" {
			return opt.Value
		}
		return opt.Name + "=" + opt.Value
	case *nodes.String:
		return opt.Str
	default:
		return ""
	}
}

// parseAGParenTuple parses a parenthesized tuple of comma-separated values,
// such as IP address tuples: ('10.0.0.1', '255.255.255.0').
// Returns the inner content formatted as comma-separated values (e.g., "'10.0.0.1', '255.255.255.0'").
func (p *Parser) parseAGParenTuple() string {
	if p.cur.Type != '(' {
		return ""
	}
	p.advance() // consume '('

	var parts []string
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}
		if p.isAGOptionToken() {
			parts = append(parts, p.agOptionTokenString())
			p.advance()
		} else {
			parts = append(parts, "UNEXPECTED:"+p.cur.Str)
			p.advance()
		}
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return strings.Join(parts, ", ")
}

// agOptionKeywords defines keywords that are valid in availability group option positions
// (both as option names and values). Core keywords like SELECT, FROM, WHERE etc. are excluded.
// Identifiers (tokIDENT) are always accepted since most AG option names (ENDPOINT_URL,
// AVAILABILITY_MODE, etc.) are unregistered as keywords.
var agOptionKeywords = map[int]bool{
	// Option names that are registered keywords
	kwROLE:        true, // ALTER AG SET (ROLE = SECONDARY)
	kwDISTRIBUTED: true, // CREATE AG WITH (DISTRIBUTED)
	kwADD:         true, // MODIFY LISTENER (...ADD IP...)
	kwREMOVE:      true, // MODIFY LISTENER (...REMOVE IP...)

	// Option values that are registered keywords
	kwNONE:      true, // DTC_SUPPORT = NONE, CLUSTER_TYPE = NONE
	kwNO:        true, // ALLOW_CONNECTIONS = NO
	kwALL:       true, // ALLOW_CONNECTIONS = ALL
	kwPRIMARY:   true, // AUTOMATED_BACKUP_PREFERENCE = PRIMARY
	kwSECONDARY: true, // AUTOMATED_BACKUP_PREFERENCE = SECONDARY
	kwMANUAL:    true, // FAILOVER_MODE = MANUAL, SEEDING_MODE = MANUAL
	kwREAD_ONLY: true, // ALLOW_CONNECTIONS = READ_ONLY
	kwEXTERNAL:  true, // FAILOVER_MODE = EXTERNAL, CLUSTER_TYPE = EXTERNAL
}

// isAGOptionToken returns true if the current token can be part of an AG option key or value.
func (p *Parser) isAGOptionToken() bool {
	if p.cur.Type == tokIDENT {
		return true
	}
	if agOptionKeywords[p.cur.Type] {
		return true
	}
	return p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST ||
		p.cur.Type == tokICONST || p.cur.Type == tokFCONST ||
		p.cur.Type == kwON || p.cur.Type == kwOFF ||
		p.cur.Type == kwNULL || p.cur.Type == kwWITH ||
		p.cur.Type == kwCREATE
}

// agOptionTokenString returns the string representation of the current token for AG options.
func (p *Parser) agOptionTokenString() string {
	switch p.cur.Type {
	case tokSCONST:
		return "'" + p.cur.Str + "'"
	case tokNSCONST:
		return "N'" + p.cur.Str + "'"
	case tokICONST, tokFCONST:
		return p.cur.Str
	default:
		return strings.ToUpper(p.cur.Str)
	}
}
