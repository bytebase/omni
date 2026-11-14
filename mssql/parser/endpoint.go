// Package parser - endpoint.go implements T-SQL ENDPOINT statement parsing.
package parser

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// SQL Server reads every word of an ENDPOINT statement by its text, so the
// sets below hold words, matched whether or not the lexer registers them as
// keywords. A bracketed or quoted word ([STARTED], [TCP]) is refused, as SQL
// Server refuses it.

// endpointTopOptionWords: endpointOption, TSql160.g:6538.
var endpointTopOptionWords = newOptionSet().withIdents("STATE", "AFFINITY")

// endpointStateValues: endpointState, TSql160.g:6544 (EndpointStateHelper).
var endpointStateValues = newOptionSet().withIdents("STARTED", "STOPPED", "DISABLED")

// endpointAffinityValues: endpointAffinity, TSql160.g:6554, takes NONE, ADMIN
// or an integer; SQL Server 2022 refuses ADMIN (Msg 7864, the dedicated
// administrator connection endpoint is reserved) and an integer (Msg 102).
var endpointAffinityValues = newOptionSet().withIdents("NONE")

// endpointProtocols: protocolInfo, TSql160.g:6588 (EndpointProtocolsHelper).
// SQL Server refuses VIA, NAMED_PIPES and SHARED_MEMORY.
var endpointProtocols = newOptionSet().withIdents("TCP", "HTTP")

// endpointTCPOptionWords and endpointHTTPOptionWords: protocolOptionsItem,
// TSql160.g:6614. SQL Server takes LISTENER_PORT and LISTENER_IP for TCP only
// and the rest for HTTP only.
var (
	endpointTCPOptionWords  = newOptionSet().withIdents("LISTENER_PORT", "LISTENER_IP")
	endpointHTTPOptionWords = newOptionSet().withIdents(
		"PATH", "SITE", "CLEAR_PORT", "SSL_PORT", "AUTHENTICATION", "PORTS",
		"AUTH_REALM", "DEFAULT_LOGON_DOMAIN", "COMPRESSION",
	)
)

// endpointHTTPAuthValues: authenticationProtocolOption, TSql160.g:6701
// (AuthenticationTypesHelper).
var endpointHTTPAuthValues = newOptionSet().withIdents("BASIC", "DIGEST", "NTLM", "KERBEROS", "INTEGRATED")

// endpointHTTPPortValues: portsProtocolOption, TSql160.g:6723 (PortTypesHelper).
var endpointHTTPPortValues = newOptionSet().withIdents("CLEAR", "SSL")

// endpointEnabledDisabled: enabledDisabled, TSql160.g:7151, and
// compressionProtocolOption, TSql160.g:6785.
var endpointEnabledDisabled = newOptionSet().withIdents("ENABLED", "DISABLED")

// endpointNoneWord: AUTH_REALM, DEFAULT_LOGON_DOMAIN and WSDL = NONE.
var endpointNoneWord = newOptionSet().withIdents("NONE")

// endpointDefaultWord: defaultLiteral in databasePayloadOption (TSql160.g:7233)
// and its siblings.
var endpointDefaultWord = newOptionSet().withIdents("DEFAULT")

// endpointPayloadTypes: payloadInfo, TSql160.g:7074 (EndpointTypesHelper, and
// DATA_MIRRORING for DATABASE_MIRRORING).
var endpointPayloadTypes = newOptionSet().withIdents(
	"TSQL", "SERVICE_BROKER", "DATABASE_MIRRORING", "DATA_MIRRORING", "SOAP",
)

// endpointPayloadOptionWords: payloadOption, TSql160.g:7112, by payload type
// (CheckIfPayloadOptionAllowed). SQL Server takes ENCRYPTION for TSQL;
// AUTHENTICATION, ENCRYPTION and the message forwarding options for
// SERVICE_BROKER; AUTHENTICATION, ENCRYPTION and ROLE for DATABASE_MIRRORING.
var endpointPayloadOptionWords = map[string]optionSet{
	"TSQL":               newOptionSet().withIdents("ENCRYPTION"),
	"SERVICE_BROKER":     newOptionSet().withIdents("AUTHENTICATION", "ENCRYPTION", "MESSAGE_FORWARDING", "MESSAGE_FORWARD_SIZE"),
	"DATABASE_MIRRORING": newOptionSet().withIdents("AUTHENTICATION", "ENCRYPTION", "ROLE"),
	"SOAP": newOptionSet().withIdents(
		"WEBMETHOD", "ADD", "ALTER", "DROP", "BATCHES", "WSDL", "SESSIONS", "LOGIN_TYPE",
		"SESSION_TIMEOUT", "DATABASE", "NAMESPACE", "SCHEMA", "CHARACTER_SET", "HEADER_LIMIT",
	),
}

// endpointAuthWords: authenticationPayloadOption, TSql160.g:7305.
var (
	endpointWindowsWord     = newOptionSet().withIdents("WINDOWS")
	endpointCertificateWord = newOptionSet().withIdents("CERTIFICATE")
	endpointAuthMethods     = newOptionSet().withIdents("NTLM", "KERBEROS", "NEGOTIATE")
)

// endpointEncryptionValues: encryptionPayloadOption, TSql160.g:7377
// (EndpointEncryptionSupportHelper). SQL Server takes every value for every
// payload type, and ALGORITHM only after SUPPORTED or REQUIRED.
var endpointEncryptionValues = newOptionSet().withIdents(
	"DISABLED", "SUPPORTED", "REQUIRED", "NEGOTIATED", "STRICT",
)

var (
	endpointAlgorithmWord  = newOptionSet().withIdents("ALGORITHM")
	endpointEncryptionAlgs = newOptionSet().withIdents("AES", "RC4")
)

// endpointRoleValues: rolePayloadOption, TSql160.g:7433.
var endpointRoleValues = newOptionSet().withIdents("WITNESS", "PARTNER", "ALL")

// SOAP payload values: TSql160.g:7164 (WSDL), 7197 (LOGIN_TYPE), 7212
// (SESSION_TIMEOUT), 7263 (SCHEMA), 7278 (CHARACTER_SET).
var (
	endpointLoginTypes     = newOptionSet().withIdents("WINDOWS", "MIXED")
	endpointNeverWord      = newOptionSet().withIdents("NEVER")
	endpointSoapSchemas    = newOptionSet().withIdents("STANDARD", "NONE")
	endpointCharacterSets  = newOptionSet().withIdents("SQL", "XML")
	endpointWebMethodWord  = newOptionSet().withIdents("WEBMETHOD")
	endpointWebMethodWords = newOptionSet().withIdents("NAME", "SCHEMA", "FORMAT")
	// soapMethodOption, TSql160.g:7504: SCHEMA = { NONE | STANDARD | DEFAULT },
	// FORMAT = { ALL_RESULTS | ROWSETS_ONLY | NONE } (SoapMethodFormatsHelper).
	endpointWebMethodSchemas = newOptionSet().withIdents("NONE", "STANDARD", "DEFAULT")
	endpointWebMethodFormats = newOptionSet().withIdents("ALL_RESULTS", "ROWSETS_ONLY", "NONE")
)

// endpointMaxInt is the largest integer SQL Server reads in an ENDPOINT
// option; a larger literal is a syntax error.
const endpointMaxInt = 2147483647

// parseCreateEndpointStmt parses CREATE ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-endpoint-transact-sql
// Grammar: createEndpointStatement, TSql160.g:6520.
//
//	CREATE ENDPOINT endPointName [ AUTHORIZATION login ]
//	    [ STATE = { STARTED | STOPPED | DISABLED } [ , AFFINITY = NONE ] ]
//	    AS { TCP | HTTP } ( <protocol_specific_arguments> )
//	    FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING | SOAP } ( <language_specific_arguments> )
//
// SQL Server also requires, while it parses, LISTENER_PORT for TCP, PATH,
// PORTS and AUTHENTICATION for HTTP, ROLE for DATABASE_MIRRORING (Msg 7801),
// and an HTTP endpoint for SOAP and a TCP one for the others (Msg 7861).
func (p *Parser) parseCreateEndpointStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// ENDPOINT keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "CREATE",
		ObjectType: "ENDPOINT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	var opts []nodes.Node
	if p.cur.Type == kwAUTHORIZATION {
		optLoc := p.pos()
		p.advance()
		login, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		opts = append(opts, &nodes.EndpointOption{Name: "AUTHORIZATION", Value: login, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	}

	ep := &endpointParse{create: true}
	if err := p.parseEndpointClauses(ep, &opts); err != nil {
		return nil, err
	}
	if ep.protocol == "" {
		return nil, p.newParseError(p.pos(), "the required parameter AS was not specified")
	}
	if ep.payload == "" {
		return nil, p.newParseError(p.pos(), "the required parameter FOR was not specified")
	}
	if (ep.protocol == "HTTP") != (ep.payload == "SOAP") {
		return nil, p.newParseError(ep.payloadLoc, fmt.Sprintf("%s endpoints cannot be of the FOR %s type", ep.protocol, ep.payload))
	}
	var required []string
	switch ep.protocol {
	case "TCP":
		required = []string{"LISTENER_PORT"}
	case "HTTP":
		required = []string{"PATH", "PORTS", "AUTHENTICATION"}
	}
	for _, r := range required {
		if !ep.protocolSeen[r] {
			return nil, p.newParseError(ep.protocolLoc, "the required parameter "+r+" was not specified")
		}
	}
	if ep.payload == "DATABASE_MIRRORING" && !ep.payloadSeen["ROLE"] {
		return nil, p.newParseError(ep.payloadLoc, "the required parameter ROLE was not specified")
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterEndpointStmt parses ALTER ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-endpoint-transact-sql
// Grammar: alterEndpointStatement, TSql160.g:11485.
//
//	ALTER ENDPOINT endPointName
//	    [ STATE = { STARTED | STOPPED | DISABLED } [ , AFFINITY = NONE ] ]
//	    [ AS { TCP | HTTP } ( <protocol_specific_arguments> ) ]
//	    [ FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING | SOAP } ( <language_specific_arguments> ) ]
//
// SQL Server refuses an ALTER ENDPOINT that changes nothing (Msg 7863).
func (p *Parser) parseAlterEndpointStmt() (*nodes.SecurityStmt, error) {
	loc := p.pos()
	// ENDPOINT keyword already consumed by caller

	stmt := &nodes.SecurityStmt{
		Action:     "ALTER",
		ObjectType: "ENDPOINT",
		Loc:        nodes.Loc{Start: loc, End: -1},
	}

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	var opts []nodes.Node
	ep := &endpointParse{}
	if err := p.parseEndpointClauses(ep, &opts); err != nil {
		return nil, err
	}
	if len(opts) == 0 {
		return nil, p.newParseError(p.pos(), "the ALTER ENDPOINT statement did not contain any values to modify")
	}

	stmt.Options = &nodes.List{Items: opts}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseDropEndpointStmt parses DROP ENDPOINT.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-endpoint-transact-sql
// Grammar: dropEndpointStatement, TSql160.g:15568.
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

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// endpointParse records what CREATE ENDPOINT checks after the clauses.
type endpointParse struct {
	create       bool
	protocol     string
	protocolLoc  int
	protocolSeen map[string]bool
	payload      string
	payloadLoc   int
	payloadSeen  map[string]bool
}

// endpointAt reports whether the current token is an undelimited word of set.
func (p *Parser) endpointAt(set optionSet) bool {
	switch {
	case p.cur.Type == tokIDENT:
		if p.isDelimited(p.cur) {
			return false
		}
	case p.cur.Type < kwABSENT:
		return false
	}
	return p.isValidOption(set)
}

// endpointExpect consumes an undelimited word of set and returns it in upper case.
func (p *Parser) endpointExpect(set optionSet) (string, error) {
	if !p.endpointAt(set) {
		return "", p.syntaxErrorAtCur()
	}
	w := strings.ToUpper(p.cur.Str)
	p.advance()
	return w, nil
}

// endpointExpectEquals consumes the '=' after an option name.
func (p *Parser) endpointExpectEquals() error {
	_, err := p.expect('=')
	return err
}

// endpointNote marks name as given in seen, refusing a second one as SQL
// Server does (Msg 7816).
func (p *Parser) endpointNote(seen map[string]bool, name string, pos int) error {
	if seen[name] {
		return p.newParseError(pos, fmt.Sprintf("a duplicate parameter was specified, '%s'", name))
	}
	seen[name] = true
	return nil
}

// endpointInteger reads an unsigned integer literal in [lo, endpointMaxInt]
// and returns its text. SQL Server reads a larger literal as a syntax error,
// and a smaller one as out of range (Msg 7810); hi narrows the range for
// ports.
func (p *Parser) endpointInteger(name string, lo, hi int64) (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	pos := p.pos()
	v := p.cur.Ival
	text := p.cur.Str
	if v > endpointMaxInt {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if v < lo || v > hi {
		return "", p.newParseError(pos, fmt.Sprintf("the value '%d' is not within range for the '%s' parameter", v, name))
	}
	return text, nil
}

// endpointString reads a non-empty string literal, as every ENDPOINT string
// option takes (nonEmptyString, TSql160.g:6745; SQL Server refuses ” for
// each), and returns it quoted.
func (p *Parser) endpointString() (string, error) {
	if (p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST) || p.cur.Str == "" {
		return "", p.syntaxErrorAtCur()
	}
	s := "'" + p.cur.Str + "'"
	p.advance()
	return s, nil
}

// endpointIsString reports whether the current token is a string literal.
func (p *Parser) endpointIsString() bool {
	return p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST
}

// parseEndpointClauses parses the clauses after the endpoint name (and
// CREATE's AUTHORIZATION), in the order SQL Server requires them:
//
//	[ endpointOption [ , endpointOption ] ] [ protocolInfo ] [ payloadInfo ]
func (p *Parser) parseEndpointClauses(ep *endpointParse, opts *[]nodes.Node) error {
	if p.endpointAt(endpointTopOptionWords) {
		seen := map[string]bool{}
		for {
			if err := p.parseEndpointTopOption(seen, opts); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	if p.cur.Type == kwAS {
		if err := p.parseEndpointProtocol(ep, opts); err != nil {
			return err
		}
	}
	if p.cur.Type == kwFOR {
		if err := p.parseEndpointPayload(ep, opts); err != nil {
			return err
		}
	}
	return nil
}

// parseEndpointTopOption parses one endpointOption, TSql160.g:6538:
//
//	STATE = { STARTED | STOPPED | DISABLED } | AFFINITY = NONE
func (p *Parser) parseEndpointTopOption(seen map[string]bool, opts *[]nodes.Node) error {
	optLoc := p.pos()
	name, err := p.endpointExpect(endpointTopOptionWords)
	if err != nil {
		return err
	}
	if err := p.endpointNote(seen, name, optLoc); err != nil {
		return err
	}
	if err := p.endpointExpectEquals(); err != nil {
		return err
	}
	values := endpointStateValues
	if name == "AFFINITY" {
		values = endpointAffinityValues
	}
	val, err := p.endpointExpect(values)
	if err != nil {
		return err
	}
	*opts = append(*opts, &nodes.EndpointOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	return nil
}

// parseEndpointProtocol parses protocolInfo, TSql160.g:6588:
//
//	AS { TCP | HTTP } ( protocolOptionsItem [ , protocolOptionsItem ]* )
func (p *Parser) parseEndpointProtocol(ep *endpointParse, opts *[]nodes.Node) error {
	optLoc := p.pos()
	p.advance() // AS
	protocol, err := p.endpointExpect(endpointProtocols)
	if err != nil {
		return err
	}
	ep.protocol = protocol
	ep.protocolLoc = optLoc
	ep.protocolSeen = map[string]bool{}
	*opts = append(*opts, &nodes.EndpointOption{Name: "AS", Value: protocol, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

	words := endpointTCPOptionWords
	if protocol == "HTTP" {
		words = endpointHTTPOptionWords
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		if err := p.parseEndpointProtocolOption(words, ep.protocolSeen, opts); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	_, err = p.expect(')')
	return err
}

// parseEndpointProtocolOption parses one protocolOptionsItem, TSql160.g:6614.
//
//	TCP:  LISTENER_PORT = port | LISTENER_IP = { ALL | ( ipv4 [ : ipv4 ] ) | ( 'ipv6' ) }
//	HTTP: PATH = 'url' | SITE = 'site' | CLEAR_PORT = port | SSL_PORT = port
//	    | AUTHENTICATION = ( { BASIC | DIGEST | NTLM | KERBEROS | INTEGRATED } [ ,...n ] )
//	    | PORTS = ( { CLEAR | SSL } [ ,...n ] )
//	    | AUTH_REALM = { 'realm' | NONE } | DEFAULT_LOGON_DOMAIN = { 'domain' | NONE }
//	    | COMPRESSION = { ENABLED | DISABLED }
func (p *Parser) parseEndpointProtocolOption(words optionSet, seen map[string]bool, opts *[]nodes.Node) error {
	optLoc := p.pos()
	name, err := p.endpointExpect(words)
	if err != nil {
		return err
	}
	if err := p.endpointNote(seen, name, optLoc); err != nil {
		return err
	}
	if err := p.endpointExpectEquals(); err != nil {
		return err
	}

	var val string
	switch name {
	case "LISTENER_PORT", "CLEAR_PORT", "SSL_PORT":
		val, err = p.endpointInteger(name, 1, 65535)
	case "LISTENER_IP":
		val, err = p.parseEndpointListenerIP()
	case "PATH", "SITE":
		val, err = p.endpointString()
	case "AUTH_REALM", "DEFAULT_LOGON_DOMAIN":
		if p.endpointAt(endpointNoneWord) {
			val = "NONE"
			p.advance()
		} else {
			val, err = p.endpointString()
		}
	case "COMPRESSION":
		val, err = p.endpointExpect(endpointEnabledDisabled)
	case "AUTHENTICATION":
		val, err = p.parseEndpointWordList(name, endpointHTTPAuthValues)
	case "PORTS":
		val, err = p.parseEndpointWordList(name, endpointHTTPPortValues)
	}
	if err != nil {
		return err
	}
	*opts = append(*opts, &nodes.EndpointOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	return nil
}

// parseEndpointWordList parses ( word [ , word ]* ) for HTTP AUTHENTICATION and
// PORTS, refusing a repeated word as SQL Server does (Msg 7816).
func (p *Parser) parseEndpointWordList(name string, values optionSet) (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	var vals []string
	seen := map[string]bool{}
	for {
		pos := p.pos()
		w, err := p.endpointExpect(values)
		if err != nil {
			return "", err
		}
		if seen[w] {
			return "", p.newParseError(pos, fmt.Sprintf("a duplicate parameter was specified, '%s'", name))
		}
		seen[w] = true
		vals = append(vals, w)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return strings.Join(vals, ","), nil
}

// parseEndpointListenerIP parses the value of LISTENER_IP,
// listenerIpProtocolOption TSql160.g:6801:
//
//	ALL | ( ipv4 [ : ipv4 ] ) | ( 'ipv6' )
//
// An IPv4 address is four decimal parts of 0 to 255 joined by dots; the lexer
// splits 10.0.75.1 into 10.0, .75 and .1. SQL Server takes only an IPv6
// address in canonical numeric form as a string (Msg 7887, 7888).
func (p *Parser) parseEndpointListenerIP() (string, error) {
	if p.cur.Type == kwALL {
		p.advance()
		return "ALL", nil
	}
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	var val string
	if p.endpointIsString() {
		pos := p.pos()
		if p.cur.Str == "" {
			return "", p.syntaxErrorAtCur()
		}
		addr, err := netip.ParseAddr(p.cur.Str)
		if err != nil || !addr.Is6() {
			return "", p.newParseError(pos, "the IPv6 address specified is not supported")
		}
		val = "'" + p.cur.Str + "'"
		p.advance()
	} else {
		first, err := p.parseEndpointIPv4()
		if err != nil {
			return "", err
		}
		val = first
		if p.cur.Type == ':' {
			p.advance()
			second, err := p.parseEndpointIPv4()
			if err != nil {
				return "", err
			}
			val += ":" + second
		}
	}
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return val, nil
}

// parseEndpointIPv4 reads the number and dot tokens of an IPv4 address,
// ipAddressV4 TSql160.g:6839. Two numbers next to each other with no dot
// between (1.5 1.5) are refused, as SQL Server refuses them.
func (p *Parser) parseEndpointIPv4() (string, error) {
	start := p.pos()
	var b strings.Builder
	prevDigit := false
	for p.cur.Type == tokICONST || p.cur.Type == tokFCONST || p.cur.Type == '.' {
		text := p.cur.Str
		if p.cur.Type == '.' {
			text = "."
		}
		if strings.ContainsAny(text, "eE") {
			return "", p.syntaxErrorAtCur()
		}
		if prevDigit && text[0] != '.' {
			return "", p.syntaxErrorAtCur()
		}
		b.WriteString(text)
		prevDigit = text[len(text)-1] != '.'
		p.advance()
	}
	addr := b.String()
	parts := strings.Split(addr, ".")
	if len(parts) != 4 {
		if addr == "" {
			return "", p.syntaxErrorAtCur()
		}
		return "", p.newParseError(start, "syntax error in IP address")
	}
	for _, part := range parts {
		if part == "" {
			return "", p.newParseError(start, "syntax error in IP address")
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return "", p.newParseError(start, "syntax error in IP address")
		}
		if n > 255 {
			return "", p.newParseError(start, fmt.Sprintf("the value '%d' is not within range for the 'IP_ADDRESS' parameter", n))
		}
	}
	return addr, nil
}

// parseEndpointPayload parses payloadInfo, TSql160.g:7074:
//
//	FOR { TSQL | SERVICE_BROKER | DATABASE_MIRRORING | DATA_MIRRORING | SOAP }
//	    ( [ payloadOption [ , payloadOption ]* ] )
func (p *Parser) parseEndpointPayload(ep *endpointParse, opts *[]nodes.Node) error {
	optLoc := p.pos()
	p.advance() // FOR
	payload, err := p.endpointExpect(endpointPayloadTypes)
	if err != nil {
		return err
	}
	if payload == "DATA_MIRRORING" {
		payload = "DATABASE_MIRRORING"
	}
	ep.payload = payload
	ep.payloadLoc = optLoc
	ep.payloadSeen = map[string]bool{}
	*opts = append(*opts, &nodes.EndpointOption{Name: "FOR", Value: payload, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})

	if _, err := p.expect('('); err != nil {
		return err
	}
	if p.cur.Type != ')' {
		for {
			if err := p.parseEndpointPayloadOption(ep, opts); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	_, err = p.expect(')')
	return err
}

// parseEndpointPayloadOption parses one payloadOption, TSql160.g:7112, of the
// options ep.payload takes.
func (p *Parser) parseEndpointPayloadOption(ep *endpointParse, opts *[]nodes.Node) error {
	optLoc := p.pos()
	name, err := p.endpointExpect(endpointPayloadOptionWords[ep.payload])
	if err != nil {
		return err
	}

	switch name {
	case "WEBMETHOD", "ADD", "ALTER", "DROP":
		return p.parseEndpointWebMethod(ep, name, optLoc, opts)
	}

	if err := p.endpointNote(ep.payloadSeen, name, optLoc); err != nil {
		return err
	}
	if err := p.endpointExpectEquals(); err != nil {
		return err
	}

	var val string
	switch name {
	case "AUTHENTICATION":
		val, err = p.parseEndpointAuthentication()
	case "ENCRYPTION":
		val, err = p.parseEndpointEncryption()
	case "ROLE":
		val, err = p.endpointExpect(endpointRoleValues)
	case "MESSAGE_FORWARDING", "BATCHES", "SESSIONS":
		val, err = p.endpointExpect(endpointEnabledDisabled)
	case "MESSAGE_FORWARD_SIZE":
		val, err = p.endpointInteger(name, 1, endpointMaxInt)
	case "HEADER_LIMIT":
		val, err = p.endpointInteger(name, 0, endpointMaxInt)
	case "SESSION_TIMEOUT":
		if p.endpointAt(endpointNeverWord) {
			val = "NEVER"
			p.advance()
		} else {
			val, err = p.endpointInteger(name, 0, endpointMaxInt)
		}
	case "WSDL":
		switch {
		case p.endpointAt(endpointNoneWord):
			val = "NONE"
			p.advance()
		case p.endpointAt(endpointDefaultWord):
			val = "DEFAULT"
			p.advance()
		default:
			val, err = p.endpointString()
		}
	case "DATABASE", "NAMESPACE":
		if p.endpointAt(endpointDefaultWord) {
			val = "DEFAULT"
			p.advance()
		} else {
			val, err = p.endpointString()
		}
	case "LOGIN_TYPE":
		val, err = p.endpointExpect(endpointLoginTypes)
	case "SCHEMA":
		val, err = p.endpointExpect(endpointSoapSchemas)
	case "CHARACTER_SET":
		val, err = p.endpointExpect(endpointCharacterSets)
	}
	if err != nil {
		return err
	}
	*opts = append(*opts, &nodes.EndpointOption{Name: name, Value: val, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	return nil
}

// parseEndpointAuthentication parses the value of a SERVICE_BROKER or
// DATABASE_MIRRORING AUTHENTICATION, authenticationPayloadOption TSql160.g:7305:
//
//	WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
//	| CERTIFICATE certificate_name
//	| WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ] CERTIFICATE certificate_name
//	| CERTIFICATE certificate_name WINDOWS [ { NTLM | KERBEROS | NEGOTIATE } ]
func (p *Parser) parseEndpointAuthentication() (string, error) {
	var parts []string
	windows := func() {
		parts = append(parts, "WINDOWS")
		p.advance()
		if p.endpointAt(endpointAuthMethods) {
			parts = append(parts, strings.ToUpper(p.cur.Str))
			p.advance()
		}
	}
	certificate := func() error {
		p.advance()
		name, ok := p.parseIdentifier()
		if !ok {
			return p.syntaxErrorAtCur()
		}
		parts = append(parts, "CERTIFICATE", name)
		return nil
	}
	switch {
	case p.endpointAt(endpointWindowsWord):
		windows()
		if p.endpointAt(endpointCertificateWord) {
			if err := certificate(); err != nil {
				return "", err
			}
		}
	case p.endpointAt(endpointCertificateWord):
		if err := certificate(); err != nil {
			return "", err
		}
		if p.endpointAt(endpointWindowsWord) {
			windows()
		}
	default:
		return "", p.syntaxErrorAtCur()
	}
	return strings.Join(parts, " "), nil
}

// parseEndpointEncryption parses the value of ENCRYPTION,
// encryptionPayloadOption TSql160.g:7377:
//
//	{ DISABLED | NEGOTIATED | STRICT
//	| { SUPPORTED | REQUIRED } [ ALGORITHM { AES | RC4 | AES RC4 | RC4 AES } ] }
func (p *Parser) parseEndpointEncryption() (string, error) {
	val, err := p.endpointExpect(endpointEncryptionValues)
	if err != nil {
		return "", err
	}
	if (val != "SUPPORTED" && val != "REQUIRED") || !p.endpointAt(endpointAlgorithmWord) {
		return val, nil
	}
	p.advance() // ALGORITHM
	alg1, err := p.endpointExpect(endpointEncryptionAlgs)
	if err != nil {
		return "", err
	}
	val += " ALGORITHM " + alg1
	if p.endpointAt(endpointEncryptionAlgs) {
		// A second algorithm must differ from the first (AES AES is refused).
		pos := p.pos()
		alg2, err := p.endpointExpect(endpointEncryptionAlgs)
		if err != nil {
			return "", err
		}
		if alg2 == alg1 {
			return "", p.newParseError(pos, "syntax error at or near \""+alg2+"\"")
		}
		val += " " + alg2
	}
	return val, nil
}

// parseEndpointWebMethod parses a SOAP web method, soapMethod TSql160.g:7457.
// CREATE ENDPOINT takes WEBMETHOD alone, ALTER ENDPOINT only ADD, ALTER and
// DROP WEBMETHOD (Msg 7803):
//
//	[ ADD | ALTER ] WEBMETHOD 'alias' ( soapMethodOption [ , soapMethodOption ]* )
//	| DROP WEBMETHOD 'alias'
//
// The first word is already consumed.
func (p *Parser) parseEndpointWebMethod(ep *endpointParse, first string, optLoc int, opts *[]nodes.Node) error {
	if first != "WEBMETHOD" {
		if _, err := p.endpointExpect(endpointWebMethodWord); err != nil {
			return err
		}
	}
	name := "WEBMETHOD"
	if first != "WEBMETHOD" {
		name = first + " WEBMETHOD"
	}
	if (first == "WEBMETHOD") != ep.create {
		stmt := "ALTER"
		if ep.create {
			stmt = "CREATE"
		}
		return p.newParseError(optLoc, fmt.Sprintf("the clause %s can not be used in the %s ENDPOINT statement", name, stmt))
	}

	// soapMethodAlias, TSql160.g:7489: 'name' | 'namespace'.'name'.
	alias, err := p.endpointString()
	if err != nil {
		return err
	}
	if p.cur.Type == '.' {
		p.advance()
		second, err := p.endpointString()
		if err != nil {
			return err
		}
		alias += "." + second
	}
	if first == "DROP" {
		*opts = append(*opts, &nodes.EndpointOption{Name: name, Value: alias, Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
		return nil
	}

	if _, err := p.expect('('); err != nil {
		return err
	}
	seen := map[string]bool{}
	var parts []string
	for {
		pos := p.pos()
		opt, err := p.endpointExpect(endpointWebMethodWords)
		if err != nil {
			return err
		}
		if err := p.endpointNote(seen, opt, pos); err != nil {
			return err
		}
		if err := p.endpointExpectEquals(); err != nil {
			return err
		}
		var v string
		switch opt {
		case "NAME":
			namePos := p.pos()
			if !p.endpointIsString() {
				return p.syntaxErrorAtCur()
			}
			if !endpointThreePartName(p.cur.Str) {
				return p.newParseError(namePos, fmt.Sprintf("the SOAP method object '%s' must be specified using a fully qualified three-part name", p.cur.Str))
			}
			v = "'" + p.cur.Str + "'"
			p.advance()
		case "SCHEMA":
			v, err = p.endpointExpect(endpointWebMethodSchemas)
		case "FORMAT":
			v, err = p.endpointExpect(endpointWebMethodFormats)
		}
		if err != nil {
			return err
		}
		parts = append(parts, opt+"="+v)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	if !seen["NAME"] {
		return p.newParseError(optLoc, "the required parameter NAME was not specified")
	}
	*opts = append(*opts, &nodes.EndpointOption{Name: name, Value: alias + " (" + strings.Join(parts, ", ") + ")", Loc: nodes.Loc{Start: optLoc, End: p.prevEnd()}})
	return nil
}

// endpointThreePartName reports whether s names an object in three non-empty
// parts, database.schema.object, where a part may be bracketed or quoted and
// hold a dot ('[a.b].c.d').
func endpointThreePartName(s string) bool {
	parts, empty := 1, true
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '.':
			if empty {
				return false
			}
			parts++
			empty = true
		case '[', '"':
			closer := byte(']')
			if c == '"' {
				closer = '"'
			}
			j := i + 1
			for ; j < len(s); j++ {
				if s[j] == closer {
					if j+1 < len(s) && s[j+1] == closer {
						j++
						continue
					}
					break
				}
			}
			if j >= len(s) {
				return false
			}
			i = j
			empty = false
		default:
			empty = false
		}
	}
	return parts == 3 && !empty
}
