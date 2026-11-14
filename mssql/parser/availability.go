// Package parser - availability.go implements T-SQL AVAILABILITY GROUP statement parsing.
//
// ScriptDOM's grammar for these statements (TSql160.g:9631-9973) predates
// most of what SQL Server 2022 accepts, so the clauses below follow the
// documentation pages cited on each function and, where they are silent, what
// the pinned SQL Server does under SET PARSEONLY ON: its parser also rejects a
// missing required replica option, a repeated option, a malformed endpoint
// URL or IP address, and an out-of-range port.
package parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// agIsWord reports whether the current token is the unquoted word word
// (upper case). A bracketed or quoted identifier never matches: SQL Server
// reads these option words only bare.
func (p *Parser) agIsWord(word string) bool {
	if kw := lookupKeyword(word); kw != tokIDENT {
		return p.cur.Type == kw
	}
	return p.cur.Type == tokIDENT && !p.isDelimited(p.cur) && strings.ToUpper(p.cur.Str) == word
}

// agWordIn returns the word of words the current token is, or "".
func (p *Parser) agWordIn(words ...string) string {
	for _, w := range words {
		if p.agIsWord(w) {
			return w
		}
	}
	return ""
}

// agExpectWord consumes the word word or fails.
func (p *Parser) agExpectWord(word string) error {
	if !p.agIsWord(word) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// agExpectWordIn consumes one of words and returns it, or fails.
func (p *Parser) agExpectWordIn(words ...string) (string, error) {
	w := p.agWordIn(words...)
	if w == "" {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return w, nil
}

// agIsString reports whether the current token is a string literal.
func (p *Parser) agIsString() bool {
	return p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST
}

// agExpectString consumes a string literal and returns its value and its
// display form ('x' or N'x').
func (p *Parser) agExpectString() (string, string, error) {
	if !p.agIsString() {
		return "", "", p.syntaxErrorAtCur()
	}
	val := p.cur.Str
	disp := "'" + val + "'"
	if p.cur.Type == tokNSCONST {
		disp = "N" + disp
	}
	p.advance()
	return val, disp, nil
}

// agExpectInt consumes an integer that fits an int (32-bit), with a leading
// minus sign when signed. SQL Server reports a larger one as out of range
// (1080) and reads anything else as a syntax error.
func (p *Parser) agExpectInt(signed bool) (int64, string, error) {
	neg := false
	if signed && p.cur.Type == '-' {
		neg = true
		p.advance()
	}
	// The lexer also makes a binary literal (0x10) a tokICONST.
	isBinary := len(p.cur.Str) > 1 && (p.cur.Str[1] == 'x' || p.cur.Str[1] == 'X')
	if p.cur.Type != tokICONST || isBinary || p.cur.Ival > 2147483647 {
		return 0, "", p.syntaxErrorAtCur()
	}
	v := p.cur.Ival
	p.advance()
	if neg {
		v = -v
	}
	return v, strconv.FormatInt(v, 10), nil
}

// agExpectName consumes a name the server parses as an identifier: a regular
// or delimited identifier, or a word that is not reserved.
func (p *Parser) agExpectName() (string, error) {
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// agExpectServerName consumes a replica or availability group name in a
// REPLICA ON or AVAILABILITY GROUP ON clause: a string literal that is not
// empty.
func (p *Parser) agExpectServerName() (string, error) {
	if !p.agIsString() || p.cur.Str == "" {
		return "", p.syntaxErrorAtCur()
	}
	name := p.cur.Str
	p.advance()
	return name, nil
}

// agValidURL reports whether s has the form SQL Server requires of an
// endpoint, listener and routing URL: TCP://system-address:port, the
// address a name, an IPv4 address or a bracketed IPv6 address, and the port
// 1 to 65535 (errors 35232, 19499, 19401, 47120).
func agValidURL(s string) bool {
	if len(s) < 6 || strings.ToUpper(s[:6]) != "TCP://" {
		return false
	}
	rest := s[6:]
	colon := strings.LastIndexByte(rest, ':')
	if colon <= 0 {
		return false
	}
	host, port := rest[:colon], rest[colon+1:]
	if host[0] == '[' {
		if host[len(host)-1] != ']' || len(host) < 3 {
			return false
		}
	} else {
		for _, c := range host {
			ok := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
				c == '.' || c == '-' || c == '_' || c == ' '
			if !ok {
				return false
			}
		}
	}
	if port == "" || len(port) > 5 {
		return false
	}
	for _, c := range port {
		if c < '0' || c > '9' {
			return false
		}
	}
	n, _ := strconv.Atoi(port)
	return n >= 1 && n <= 65535
}

// agIPv4 and agIPv6 check a listener address (error 19483).
func agIPv4(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && !strings.Contains(s, ":") && ip.To4() != nil
}

func agIPv6(s string) bool {
	return strings.Contains(s, ":") && net.ParseIP(s) != nil
}

// agError reports a semantic error at pos.
func (p *Parser) agError(pos int, format string, args ...any) error {
	return p.newParseError(pos, fmt.Sprintf(format, args...))
}

// parseCreateAvailabilityGroupStmt parses CREATE AVAILABILITY GROUP.
// Caller has consumed CREATE AVAILABILITY GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-availability-group-transact-sql
//
//	CREATE AVAILABILITY GROUP group_name
//	   [ WITH ( <with_option_spec> [ ,...n ] ) ]
//	   { FOR [ DATABASE database_name [ ,...n ] ] REPLICA ON <add_replica_spec> [ ,...n ]
//	       [ LISTENER 'dns_name' ( <listener_option> ) ]
//	   | AVAILABILITY GROUP ON <add_availability_group_spec> [ ,...n ] }
//
// AVAILABILITY GROUP ON needs WITH (DISTRIBUTED) (error 19498), and
// DISTRIBUTED takes no other option beside it (19497).
func (p *Parser) parseCreateAvailabilityGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.agExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	distributed := false
	if p.cur.Type == kwWITH {
		start := p.pos()
		p.advance()
		opts = append(opts, newAGOpt("WITH", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		groupOpts, dist, err := p.parseAGGroupOptionList(agCreateGroupOptions)
		if err != nil {
			return nil, err
		}
		distributed = dist
		opts = append(opts, groupOpts...)
	}

	switch {
	case p.cur.Type == kwFOR:
		start := p.pos()
		p.advance()
		if p.cur.Type == kwDATABASE {
			p.advance()
			var dbs []string
			seen := map[string]bool{}
			for {
				dbPos := p.pos()
				db, err := p.agExpectName()
				if err != nil {
					return nil, err
				}
				if seen[strings.ToUpper(db)] {
					return nil, p.agError(dbPos, "database name '%s' was specified more than once", db)
				}
				seen[strings.ToUpper(db)] = true
				dbs = append(dbs, db)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			opts = append(opts, newAGOpt("FOR DATABASE", strings.Join(dbs, ", "), nodes.Loc{Start: start, End: p.prevEnd()}))
		} else {
			opts = append(opts, newAGOpt("FOR", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		}
		start = p.pos()
		if _, err := p.expect(kwREPLICA); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		opts = append(opts, newAGOpt("REPLICA ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		specs, err := p.parseAGReplicaSpecs(agReplicaAdd, true)
		if err != nil {
			return nil, err
		}
		opts = append(opts, specs...)
		if p.cur.Type == kwLISTENER {
			start := p.pos()
			p.advance()
			disp, err := p.parseAGListenerName()
			if err != nil {
				return nil, err
			}
			opts = append(opts, newAGOpt("LISTENER", disp, nodes.Loc{Start: start, End: p.prevEnd()}))
			lopts, err := p.parseAGListenerOptions()
			if err != nil {
				return nil, err
			}
			opts = append(opts, lopts...)
		}

	case p.cur.Type == kwAVAILABILITY:
		start := p.pos()
		if !distributed {
			return nil, p.agError(start, "to create a distributed availability group, specify the DISTRIBUTED option")
		}
		p.advance()
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwON); err != nil {
			return nil, err
		}
		opts = append(opts, newAGOpt("AVAILABILITY GROUP ON", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		specs, err := p.parseAGReplicaSpecs(agGroupAdd, true)
		if err != nil {
			return nil, err
		}
		opts = append(opts, specs...)

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterAvailabilityGroupStmt parses ALTER AVAILABILITY GROUP, which
// takes exactly one action.
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
//	   | JOIN [ WITH ( CLUSTER_TYPE = { WSFC | EXTERNAL | NONE } ) ]
//	   | JOIN AVAILABILITY GROUP ON <add_availability_group_spec> [ ,...n ]
//	   | MODIFY AVAILABILITY GROUP ON <modify_availability_group_spec> [ ,...n ]
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
//
// SQL Server takes one database, one replica and one SET option per
// statement, and a MODIFY REPLICA that changes one setting (error 10786).
func (p *Parser) parseAlterAvailabilityGroupStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.agExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	var opts []nodes.Node
	start := p.pos()
	mark := func(name, value string) {
		opts = append(opts, newAGOpt(name, value, nodes.Loc{Start: start, End: p.prevEnd()}))
	}
	switch {
	case p.cur.Type == kwSET:
		p.advance()
		mark("SET", "")
		setOpts, _, err := p.parseAGGroupOptionList(agSetGroupOptions)
		if err != nil {
			return nil, err
		}
		opts = append(opts, setOpts...)

	case p.cur.Type == kwADD || p.cur.Type == kwREMOVE || p.cur.Type == kwMODIFY:
		verb := strings.ToUpper(p.cur.Str)
		p.advance()
		switch {
		case p.cur.Type == kwDATABASE && verb != "MODIFY":
			p.advance()
			db, err := p.agExpectName()
			if err != nil {
				return nil, err
			}
			mark(verb+" DATABASE", db)
		case p.cur.Type == kwREPLICA:
			p.advance()
			if _, err := p.expect(kwON); err != nil {
				return nil, err
			}
			mark(verb+" REPLICA ON", "")
			switch verb {
			case "ADD":
				specs, err := p.parseAGReplicaSpecs(agReplicaAdd, false)
				if err != nil {
					return nil, err
				}
				opts = append(opts, specs...)
			case "MODIFY":
				specs, err := p.parseAGReplicaSpecs(agReplicaModify, false)
				if err != nil {
					return nil, err
				}
				opts = append(opts, specs...)
			default:
				namePos := p.pos()
				server, err := p.agExpectServerName()
				if err != nil {
					return nil, err
				}
				opts = append(opts, newAGOpt("'"+server+"'", "", nodes.Loc{Start: namePos, End: p.prevEnd()}))
			}
		case p.cur.Type == kwLISTENER:
			p.advance()
			disp, err := p.parseAGListenerName()
			if err != nil {
				return nil, err
			}
			mark(verb+" LISTENER", disp)
			switch verb {
			case "ADD":
				lopts, err := p.parseAGListenerOptions()
				if err != nil {
					return nil, err
				}
				opts = append(opts, lopts...)
			case "MODIFY":
				lopts, err := p.parseAGModifyListenerOption()
				if err != nil {
					return nil, err
				}
				opts = append(opts, lopts...)
			}
		case p.cur.Type == kwAVAILABILITY && verb == "MODIFY":
			p.advance()
			if _, err := p.expect(kwGROUP); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwON); err != nil {
				return nil, err
			}
			mark("MODIFY AVAILABILITY GROUP ON", "")
			specs, err := p.parseAGReplicaSpecs(agGroupModify, true)
			if err != nil {
				return nil, err
			}
			opts = append(opts, specs...)
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case p.cur.Type == kwRESTART:
		p.advance()
		if _, err := p.expect(kwLISTENER); err != nil {
			return nil, err
		}
		disp, err := p.parseAGListenerName()
		if err != nil {
			return nil, err
		}
		mark("RESTART LISTENER", disp)

	case p.cur.Type == kwJOIN:
		p.advance()
		switch p.cur.Type {
		case kwAVAILABILITY:
			p.advance()
			if _, err := p.expect(kwGROUP); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwON); err != nil {
				return nil, err
			}
			mark("JOIN AVAILABILITY GROUP ON", "")
			specs, err := p.parseAGReplicaSpecs(agGroupAdd, true)
			if err != nil {
				return nil, err
			}
			opts = append(opts, specs...)
		case kwWITH:
			mark("JOIN", "")
			withPos := p.pos()
			p.advance()
			opts = append(opts, newAGOpt("WITH", "", nodes.Loc{Start: withPos, End: p.prevEnd()}))
			joinOpts, _, err := p.parseAGGroupOptionList(agJoinGroupOptions)
			if err != nil {
				return nil, err
			}
			opts = append(opts, joinOpts...)
		default:
			mark("JOIN", "")
		}

	case p.cur.Type == kwGRANT || p.cur.Type == kwDENY:
		verb := strings.ToUpper(p.cur.Str)
		p.advance()
		for _, kw := range []int{kwCREATE, kwANY, kwDATABASE} {
			if _, err := p.expect(kw); err != nil {
				return nil, err
			}
		}
		mark(verb+" CREATE ANY DATABASE", "")

	case p.cur.Type == kwFAILOVER:
		p.advance()
		mark("FAILOVER", "")

	case p.cur.Type == kwFORCE_FAILOVER_ALLOW_DATA_LOSS:
		p.advance()
		mark("FORCE_FAILOVER_ALLOW_DATA_LOSS", "")

	case p.cur.Type == kwOFFLINE:
		p.advance()
		mark("OFFLINE", "")

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Options = &nodes.List{Items: opts}
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
	stmt := &nodes.SecurityStmt{
		Action:     "DROP",
		ObjectType: "AVAILABILITY GROUP",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}
	name, err := p.agExpectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// newAGOpt creates an AvailabilityGroupOption with the given name, value, and location.
func newAGOpt(name, value string, loc nodes.Loc) *nodes.AvailabilityGroupOption {
	return &nodes.AvailabilityGroupOption{Name: name, Value: value, Loc: loc}
}

// agGroupOptionContext names where an availability group option list is.
type agGroupOptionContext int

const (
	agCreateGroupOptions agGroupOptionContext = iota // CREATE ... WITH ( ... )
	agSetGroupOptions                                // ALTER ... SET ( one option )
	agJoinGroupOptions                               // ALTER ... JOIN WITH ( ... )
)

// parseAGGroupOptionList parses ( <with_option_spec> [ ,...n ] ) and reports
// whether DISTRIBUTED was among them.
//
//	<with_option_spec> ::=
//	    AUTOMATED_BACKUP_PREFERENCE = { PRIMARY | SECONDARY_ONLY | SECONDARY | NONE }
//	  | FAILURE_CONDITION_LEVEL = integer
//	  | HEALTH_CHECK_TIMEOUT = milliseconds
//	  | DB_FAILOVER = { ON | OFF }
//	  | DTC_SUPPORT = { PER_DB | NONE }
//	  | [ BASIC | DISTRIBUTED | CONTAINED [ REUSE_SYSTEM_DATABASES | AUTOSEEDING_SYSTEM_DATABASES ] ]
//	  | REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT = integer
//	  | CLUSTER_TYPE = { WSFC | EXTERNAL | NONE }
//	  | WRITE_LEASE_VALIDITY = seconds
//
// CREATE takes all of them but no ROLE; SET takes one, and not CLUSTER_TYPE,
// CONTAINED or the system database words, but ROLE = { PRIMARY | SECONDARY };
// JOIN WITH takes CLUSTER_TYPE only. SQL Server refuses a repeated option
// (1039), except DISTRIBUTED, and DISTRIBUTED beside any other (19497).
// CLUSTER_CONNECTION_OPTIONS is SQL Server 2025's, and 2022 refuses it.
func (p *Parser) parseAGGroupOptionList(ctx agGroupOptionContext) ([]nodes.Node, bool, error) {
	if _, err := p.expect('('); err != nil {
		return nil, false, err
	}
	var names []string
	switch ctx {
	case agCreateGroupOptions:
		names = []string{"AUTOMATED_BACKUP_PREFERENCE", "FAILURE_CONDITION_LEVEL", "HEALTH_CHECK_TIMEOUT",
			"DB_FAILOVER", "DTC_SUPPORT", "REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT", "CLUSTER_TYPE",
			"WRITE_LEASE_VALIDITY", "BASIC", "DISTRIBUTED", "CONTAINED", "REUSE_SYSTEM_DATABASES",
			"AUTOSEEDING_SYSTEM_DATABASES"}
	case agSetGroupOptions:
		names = []string{"AUTOMATED_BACKUP_PREFERENCE", "FAILURE_CONDITION_LEVEL", "HEALTH_CHECK_TIMEOUT",
			"DB_FAILOVER", "DTC_SUPPORT", "REQUIRED_SYNCHRONIZED_SECONDARIES_TO_COMMIT", "ROLE",
			"WRITE_LEASE_VALIDITY", "BASIC", "DISTRIBUTED"}
	case agJoinGroupOptions:
		names = []string{"CLUSTER_TYPE"}
	}
	var opts []nodes.Node
	seen := map[string]bool{}
	distributed, others := false, false
	for {
		start := p.pos()
		name, err := p.agExpectWordIn(names...)
		if err != nil {
			return nil, false, err
		}
		if seen[name] && name != "DISTRIBUTED" {
			return nil, false, p.agError(start, "option '%s' is specified more than once", name)
		}
		seen[name] = true
		value := ""
		switch name {
		case "BASIC", "DISTRIBUTED", "CONTAINED", "REUSE_SYSTEM_DATABASES", "AUTOSEEDING_SYSTEM_DATABASES":
		default:
			if _, err := p.expect('='); err != nil {
				return nil, false, err
			}
			switch name {
			case "AUTOMATED_BACKUP_PREFERENCE":
				value, err = p.agExpectWordIn("PRIMARY", "SECONDARY_ONLY", "SECONDARY", "NONE")
			case "DB_FAILOVER":
				value, err = p.agExpectWordIn("ON", "OFF")
			case "DTC_SUPPORT":
				value, err = p.agExpectWordIn("PER_DB", "NONE")
			case "CLUSTER_TYPE":
				value, err = p.agExpectWordIn("WSFC", "EXTERNAL", "NONE")
			case "ROLE":
				value, err = p.agExpectWordIn("SECONDARY", "PRIMARY")
			default:
				_, value, err = p.agExpectInt(false)
			}
			if err != nil {
				return nil, false, err
			}
		}
		if name == "DISTRIBUTED" {
			distributed = true
		} else {
			others = true
		}
		if distributed && others {
			return nil, false, p.agError(start, "the DISTRIBUTED option cannot be combined with any other availability group option")
		}
		opts = append(opts, newAGOpt(name, value, nodes.Loc{Start: start, End: p.prevEnd()}))
		if p.cur.Type != ',' || ctx == agSetGroupOptions {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, false, err
	}
	return opts, distributed, nil
}

// agReplicaKind names which replica spec a REPLICA ON or AVAILABILITY GROUP
// ON clause holds.
type agReplicaKind int

const (
	agReplicaAdd    agReplicaKind = iota // CREATE ... REPLICA ON, ADD REPLICA ON
	agReplicaModify                      // MODIFY REPLICA ON
	agGroupAdd                           // CREATE ... AVAILABILITY GROUP ON, JOIN AVAILABILITY GROUP ON
	agGroupModify                        // MODIFY AVAILABILITY GROUP ON
)

// parseAGReplicaSpecs parses 'name' WITH ( options ) [ , ...n ] (one spec
// unless many) and returns each spec as its name, a WITH marker and its
// options. SQL Server refuses a name given twice (35235), compared without
// case.
func (p *Parser) parseAGReplicaSpecs(kind agReplicaKind, many bool) ([]nodes.Node, error) {
	var opts []nodes.Node
	seen := map[string]bool{}
	for {
		start := p.pos()
		name, err := p.agExpectServerName()
		if err != nil {
			return nil, err
		}
		if seen[strings.ToUpper(name)] {
			return nil, p.agError(start, "the system name '%s' was specified more than once", name)
		}
		seen[strings.ToUpper(name)] = true
		opts = append(opts, newAGOpt("'"+name+"'", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		withPos := p.pos()
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		opts = append(opts, newAGOpt("WITH", "", nodes.Loc{Start: withPos, End: p.prevEnd()}))
		specOpts, err := p.parseAGReplicaOptions(kind, start)
		if err != nil {
			return nil, err
		}
		opts = append(opts, specOpts...)
		if !many || p.cur.Type != ',' {
			return opts, nil
		}
		p.advance()
	}
}

// parseAGReplicaOptions parses the parenthesized options of one replica or
// availability group spec.
//
//	<add_replica_spec> ::= <server_instance> WITH (
//	       ENDPOINT_URL = 'TCP://system-address:port',
//	       AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT | CONFIGURATION_ONLY },
//	       FAILOVER_MODE = { AUTOMATIC | MANUAL | EXTERNAL }
//	       [ , <add_replica_option> [ ,...n ] ] )
//	<add_replica_option> ::=
//	       SEEDING_MODE = { AUTOMATIC | MANUAL }
//	     | BACKUP_PRIORITY = n
//	     | SECONDARY_ROLE ( ... )
//	     | PRIMARY_ROLE ( ... )
//	     | SESSION_TIMEOUT = integer
//	<add_availability_group_spec> ::= <ag_name> WITH (
//	       LISTENER_URL = 'TCP://system-address:port',
//	       AVAILABILITY_MODE = { SYNCHRONOUS_COMMIT | ASYNCHRONOUS_COMMIT },
//	       FAILOVER_MODE = MANUAL,
//	       SEEDING_MODE = { AUTOMATIC | MANUAL } )
//
// The options come in any order. Adding a replica needs ENDPOINT_URL (35236),
// AVAILABILITY_MODE (41198) and FAILOVER_MODE (41197), and a
// CONFIGURATION_ONLY one takes no FAILOVER_MODE, SEEDING_MODE,
// BACKUP_PRIORITY or role (47110); adding an availability group needs
// LISTENER_URL and AVAILABILITY_MODE and fails over only manually (35277).
// MODIFY REPLICA changes one setting (10786). SESSION_TIMEOUT must be
// positive (153).
func (p *Parser) parseAGReplicaOptions(kind agReplicaKind, specPos int) ([]nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	group := kind == agGroupAdd || kind == agGroupModify
	names := []string{"ENDPOINT_URL", "AVAILABILITY_MODE", "FAILOVER_MODE", "SEEDING_MODE",
		"BACKUP_PRIORITY", "SESSION_TIMEOUT", "SECONDARY_ROLE", "PRIMARY_ROLE"}
	if group {
		names = []string{"LISTENER_URL", "AVAILABILITY_MODE", "FAILOVER_MODE", "SEEDING_MODE", "SESSION_TIMEOUT"}
	}
	var opts []nodes.Node
	seen := map[string]int{}
	configOnly := false
	for {
		start := p.pos()
		name, err := p.agExpectWordIn(names...)
		if err != nil {
			return nil, err
		}
		if _, dup := seen[name]; dup {
			return nil, p.agError(start, "option '%s' is specified more than once", name)
		}
		seen[name] = start
		if kind == agReplicaModify && len(seen) > 1 {
			return nil, p.agError(start, "MODIFY REPLICA changes one replica option at a time")
		}
		var opt *nodes.AvailabilityGroupOption
		switch name {
		case "SECONDARY_ROLE", "PRIMARY_ROLE":
			opt, err = p.parseAGRoleOption(name, start, kind == agReplicaModify)
			if err != nil {
				return nil, err
			}
		default:
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			valPos := p.pos()
			var value string
			switch name {
			case "ENDPOINT_URL", "LISTENER_URL":
				var url string
				url, value, err = p.agExpectString()
				if err == nil && !agValidURL(url) {
					err = p.agError(valPos, "the specified %s '%s' is invalid", name, url)
				}
			case "AVAILABILITY_MODE":
				value, err = p.agExpectWordIn("SYNCHRONOUS_COMMIT", "ASYNCHRONOUS_COMMIT", "CONFIGURATION_ONLY")
				configOnly = value == "CONFIGURATION_ONLY"
			case "FAILOVER_MODE":
				value, err = p.agExpectWordIn("AUTOMATIC", "MANUAL", "EXTERNAL")
				if err == nil && group && value != "MANUAL" {
					err = p.agError(valPos, "automatic failover is not supported for distributed availability group replica")
				}
			case "SEEDING_MODE":
				value, err = p.agExpectWordIn("AUTOMATIC", "MANUAL")
			case "BACKUP_PRIORITY":
				_, value, err = p.agExpectInt(true)
			case "SESSION_TIMEOUT":
				var n int64
				n, value, err = p.agExpectInt(true)
				if err == nil && n <= 0 {
					err = p.agError(valPos, "invalid usage of the option SESSION_TIMEOUT")
				}
			}
			if err != nil {
				return nil, err
			}
			opt = newAGOpt(name, value, nodes.Loc{Start: start, End: p.prevEnd()})
		}
		opts = append(opts, opt)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	switch kind {
	case agReplicaAdd:
		for _, req := range []string{"ENDPOINT_URL", "AVAILABILITY_MODE"} {
			if _, ok := seen[req]; !ok {
				return nil, p.agError(specPos, "the %s option has not been specified for the replica", req)
			}
		}
		if configOnly {
			for _, bad := range []string{"FAILOVER_MODE", "SEEDING_MODE", "BACKUP_PRIORITY", "SECONDARY_ROLE", "PRIMARY_ROLE"} {
				if pos, ok := seen[bad]; ok {
					return nil, p.agError(pos, "the %s option is not valid for a configuration-only replica", bad)
				}
			}
		} else if _, ok := seen["FAILOVER_MODE"]; !ok {
			return nil, p.agError(specPos, "the FAILOVER_MODE option has not been specified for the replica")
		}
	case agGroupAdd:
		for _, req := range []string{"LISTENER_URL", "AVAILABILITY_MODE"} {
			if _, ok := seen[req]; !ok {
				return nil, p.agError(specPos, "the %s option has not been specified for the availability group", req)
			}
		}
	}
	return opts, nil
}

// parseAGRoleOption parses SECONDARY_ROLE ( ... ) or PRIMARY_ROLE ( ... );
// the role word, which starts at start, is consumed. Each option inside appears at most once (1039),
// and in MODIFY REPLICA only one (10786).
//
//	SECONDARY_ROLE ( [ ALLOW_CONNECTIONS = { NO | READ_ONLY | ALL } ]
//	    [,] [ READ_ONLY_ROUTING_URL = 'TCP://system-address:port' ] )
//	PRIMARY_ROLE ( [ ALLOW_CONNECTIONS = { READ_WRITE | ALL } ]
//	    [,] [ READ_ONLY_ROUTING_LIST = { ( '<server_instance>' [ ,...n ] ) | NONE } ]
//	    [,] [ READ_WRITE_ROUTING_URL = 'TCP://system-address:port' ] )
//
// The routing list may group instances for load balancing:
// ( ( 'a', 'b' ), 'c' ), and names no instance twice (19402).
func (p *Parser) parseAGRoleOption(role string, start int, single bool) (*nodes.AvailabilityGroupOption, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	names, allow := []string{"ALLOW_CONNECTIONS", "READ_ONLY_ROUTING_URL"}, []string{"NO", "READ_ONLY", "ALL"}
	if role == "PRIMARY_ROLE" {
		names = []string{"ALLOW_CONNECTIONS", "READ_ONLY_ROUTING_LIST", "READ_WRITE_ROUTING_URL"}
		allow = []string{"READ_WRITE", "ALL"}
	}
	var parts []string
	seen := map[string]bool{}
	for {
		optPos := p.pos()
		name, err := p.agExpectWordIn(names...)
		if err != nil {
			return nil, err
		}
		if seen[name] {
			return nil, p.agError(optPos, "option '%s' is specified more than once", name)
		}
		seen[name] = true
		if single && len(seen) > 1 {
			return nil, p.agError(optPos, "MODIFY REPLICA changes one replica option at a time")
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valPos := p.pos()
		var value string
		switch name {
		case "ALLOW_CONNECTIONS":
			value, err = p.agExpectWordIn(allow...)
		case "READ_ONLY_ROUTING_LIST":
			value, err = p.parseAGRoutingList()
		default:
			var url string
			url, value, err = p.agExpectString()
			if err == nil && !agValidURL(url) {
				err = p.agError(valPos, "the %s '%s' is not valid", name, url)
			}
		}
		if err != nil {
			return nil, err
		}
		parts = append(parts, name+"="+value)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return newAGOpt(role+"("+strings.Join(parts, ", ")+")", "", nodes.Loc{Start: start, End: p.prevEnd()}), nil
}

// parseAGRoutingList parses NONE or ( item [ ,...n ] ), each item a string
// or a parenthesized list of strings, and returns its display form.
func (p *Parser) parseAGRoutingList() (string, error) {
	if p.agIsWord("NONE") {
		p.advance()
		return "NONE", nil
	}
	seen := map[string]bool{}
	instance := func() (string, error) {
		pos := p.pos()
		val, disp, err := p.agExpectString()
		if err != nil {
			return "", err
		}
		if seen[strings.ToUpper(val)] {
			return "", p.agError(pos, "a duplicate availability replica '%s' was specified in the READ_ONLY_ROUTING_LIST", val)
		}
		seen[strings.ToUpper(val)] = true
		return disp, nil
	}
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	var items []string
	for {
		if p.cur.Type == '(' {
			p.advance()
			var group []string
			for {
				disp, err := instance()
				if err != nil {
					return "", err
				}
				group = append(group, disp)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
			items = append(items, "("+strings.Join(group, ", ")+")")
		} else {
			disp, err := instance()
			if err != nil {
				return "", err
			}
			items = append(items, disp)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return "(" + strings.Join(items, ", ") + ")", nil
}

// parseAGListenerName consumes a listener's DNS name, a string of 1 to 63
// characters (19485), and returns its display form.
func (p *Parser) parseAGListenerName() (string, error) {
	pos := p.pos()
	name, _, err := p.agExpectString()
	if err != nil {
		return "", err
	}
	if n := len([]rune(name)); n < 1 || n > 63 {
		return "", p.agError(pos, "the DNS name '%s' is invalid: it must be 1 to 63 characters", name)
	}
	return "'" + name + "'", nil
}

// parseAGListenerOptions parses a listener's ( <listener_option> ).
//
//	<listener_option> ::=
//	    WITH DHCP [ ON ( 'four_part_ipv4_address', 'four_part_ipv4_mask' ) ]
//	  | WITH IP ( { ( 'four_part_ipv4_address', 'four_part_ipv4_mask' ) | ( 'ipv6_address' ) } [ ,...n ] )
//	  [ , PORT = listener_port ]
//
// SQL Server checks each address's form (19483), refuses an address listed
// twice (19484) and a port outside 1 to 65535 (19482).
func (p *Parser) parseAGListenerOptions() ([]nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	start := p.pos()
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	opts = append(opts, newAGOpt("WITH", "", nodes.Loc{Start: start, End: p.prevEnd()}))
	start = p.pos()
	kind, err := p.agExpectWordIn("DHCP", "IP")
	if err != nil {
		return nil, err
	}
	if kind == "DHCP" {
		opts = append(opts, newAGOpt("DHCP", "", nodes.Loc{Start: start, End: p.prevEnd()}))
		if p.cur.Type == kwON {
			start = p.pos()
			p.advance()
			tuple, err := p.parseAGIPTuple(false, nil)
			if err != nil {
				return nil, err
			}
			opts = append(opts, newAGOpt("ON"+tuple, "", nodes.Loc{Start: start, End: p.prevEnd()}))
		}
	} else {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		seen := map[string]bool{}
		var tuples []string
		for {
			tuple, err := p.parseAGIPTuple(true, seen)
			if err != nil {
				return nil, err
			}
			tuples = append(tuples, tuple)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		opts = append(opts, newAGOpt("IP("+strings.Join(tuples, ", ")+")", "", nodes.Loc{Start: start, End: p.prevEnd()}))
	}
	if p.cur.Type == ',' {
		p.advance()
		port, err := p.parseAGPort()
		if err != nil {
			return nil, err
		}
		opts = append(opts, port)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseAGModifyListenerOption parses MODIFY LISTENER's single option.
//
//	<modify_listener_option> ::=
//	    ADD IP ( { 'four_part_ipv4_address', 'four_part_ipv4_mask' | 'ipv6_address' } )
//	  | PORT = listener_port
//	  | REMOVE IP ( { 'four_part_ipv4_address' | 'ipv6_address' } )
func (p *Parser) parseAGModifyListenerOption() ([]nodes.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var opts []nodes.Node
	start := p.pos()
	switch {
	case p.cur.Type == kwADD || p.cur.Type == kwREMOVE:
		verb := strings.ToUpper(p.cur.Str)
		p.advance()
		opts = append(opts, newAGOpt(verb, "", nodes.Loc{Start: start, End: p.prevEnd()}))
		start = p.pos()
		if err := p.agExpectWord("IP"); err != nil {
			return nil, err
		}
		var tuple string
		if verb == "ADD" {
			t, err := p.parseAGIPTuple(true, nil)
			if err != nil {
				return nil, err
			}
			tuple = t
		} else {
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			pos := p.pos()
			addr, disp, err := p.agExpectString()
			if err != nil {
				return nil, err
			}
			if !agIPv4(addr) && !agIPv6(addr) {
				return nil, p.agError(pos, "the format for the IP address '%s' is invalid", addr)
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			tuple = "(" + disp + ")"
		}
		opts = append(opts, newAGOpt("IP"+tuple, "", nodes.Loc{Start: start, End: p.prevEnd()}))
	default:
		port, err := p.parseAGPort()
		if err != nil {
			return nil, err
		}
		opts = append(opts, port)
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseAGIPTuple parses ( 'ipv4', 'ipv4_mask' ), or ( 'ipv6' ) when v6 is
// allowed, and returns its display form. seen, when not nil, collects the
// addresses of an IP list to refuse one listed twice.
func (p *Parser) parseAGIPTuple(v6 bool, seen map[string]bool) (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	pos := p.pos()
	addr, disp, err := p.agExpectString()
	if err != nil {
		return "", err
	}
	parts := []string{disp}
	switch {
	case agIPv4(addr):
		if _, err := p.expect(','); err != nil {
			return "", err
		}
		maskPos := p.pos()
		mask, maskDisp, err := p.agExpectString()
		if err != nil {
			return "", err
		}
		if !agIPv4(mask) {
			return "", p.agError(maskPos, "the format for the IP address '%s' is invalid", mask)
		}
		parts = append(parts, maskDisp)
	case v6 && agIPv6(addr):
	default:
		return "", p.agError(pos, "the format for the IP address '%s' is invalid", addr)
	}
	if seen != nil {
		key := net.ParseIP(addr).String()
		if seen[key] {
			return "", p.agError(pos, "the IP address '%s' is duplicated in the IP address list", addr)
		}
		seen[key] = true
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return "(" + strings.Join(parts, ", ") + ")", nil
}

// parseAGPort parses PORT = listener_port, a port from 1 to 65535 (19482).
func (p *Parser) parseAGPort() (*nodes.AvailabilityGroupOption, error) {
	start := p.pos()
	if err := p.agExpectWord("PORT"); err != nil {
		return nil, err
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	valPos := p.pos()
	n, value, err := p.agExpectInt(false)
	if err != nil {
		return nil, err
	}
	if n < 1 || n > 65535 {
		return nil, p.agError(valPos, "the port number %d is invalid: valid port numbers range from 1 to 65535", n)
	}
	return newAGOpt("PORT", value, nodes.Loc{Start: start, End: p.prevEnd()}), nil
}
