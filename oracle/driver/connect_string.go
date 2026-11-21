package driver

import (
	"regexp"
	"strings"
)

// Endpoint is one listener address an Oracle connect string may reach.
type Endpoint struct {
	Host string
	Port string
}

// ConnectTarget is what an Oracle connect string names: one service or SID,
// reachable at one or more listener addresses. Every address in a single
// descriptor is treated as an alternative for the same database, as in RAC
// address lists and Data Guard description lists.
type ConnectTarget struct {
	Endpoints []Endpoint
	// Exactly one of Service and SID is set, as written. They name different
	// things on a listener and should not be matched against each other.
	Service string
	SID     string
}

const defaultConnectPort = "1521"

var (
	tnsKeywordPattern = regexp.MustCompile(`(?i)\(\s*(HOST|PORT|SERVICE_NAME|SID)\s*=\s*([^()]+?)\s*\)`)
	tnsAddressPattern = regexp.MustCompile(`(?i)\(\s*ADDRESS\s*=`)
	// tnsSourceRoutePattern matches a descriptor whose addresses are a hop chain
	// through Connection Manager rather than alternatives.
	tnsSourceRoutePattern = regexp.MustCompile(`(?i)\(\s*SOURCE_ROUTE\s*=\s*(YES|ON|TRUE)\s*\)`)
	// connectSchemePattern matches the Easy Connect Plus `tcp://` / `tcps://` prefix.
	connectSchemePattern = regexp.MustCompile(`(?i)^tcps?://`)
)

// ParseConnectString parses an Oracle connect string into its listener
// endpoints and the single service or SID it names.
//
// Accepted forms, per the Oracle Net Services Administrator's Guide:
//   - Easy Connect / Easy Connect Plus:
//     [tcp[s]://][//]host1[,host2][:port1][,host3:port2]/service[:server][/instance][?params]
//     where a port applies to every host listed since the previous port, and
//     hosts without a port default to 1521.
//   - TNS connect descriptors whose ADDRESSes are alternatives for one service
//     or SID (DESCRIPTION, DESCRIPTION_LIST, ADDRESS_LIST).
//
// It returns ok=false for anything that does not name a single database: a
// bare TNS alias, a descriptor with no host or with multiple distinct services
// or a mixed service/SID, or a SOURCE_ROUTE hop chain whose first address is
// a proxy rather than the database.
func ParseConnectString(connectString string) (ConnectTarget, bool) {
	s := strings.TrimSpace(connectString)
	if s == "" {
		return ConnectTarget{}, false
	}
	if strings.HasPrefix(s, "(") {
		return parseTNSDescriptor(s)
	}
	return parseEasyConnect(s)
}

func parseEasyConnect(s string) (ConnectTarget, bool) {
	if loc := connectSchemePattern.FindStringIndex(s); loc != nil {
		s = s[loc[1]:]
	} else {
		s = strings.TrimPrefix(s, "//")
	}
	hostPart, rest, found := strings.Cut(s, "/")
	if !found || hostPart == "" {
		return ConnectTarget{}, false
	}
	service := rest
	if i := strings.IndexAny(service, ":/?"); i >= 0 {
		service = service[:i]
	}
	service = strings.TrimSpace(service)
	if service == "" || strings.ContainsAny(hostPart, " \t") {
		return ConnectTarget{}, false
	}
	endpoints, ok := parseEasyConnectHosts(hostPart)
	if !ok {
		return ConnectTarget{}, false
	}
	// Easy Connect names a service; a SID requires a TNS descriptor.
	return ConnectTarget{Endpoints: endpoints, Service: service}, true
}

// parseEasyConnectHosts splits `host1[,host2][:port1][,host3:port2]`. A port
// belongs to every host listed since the previous port; hosts with no port
// after them get Oracle's default.
func parseEasyConnectHosts(hostPart string) ([]Endpoint, bool) {
	var endpoints []Endpoint
	unported := 0
	for _, item := range strings.Split(hostPart, ",") {
		host, port, hasPort, ok := splitHostPort(item)
		if !ok {
			return nil, false
		}
		endpoints = append(endpoints, Endpoint{Host: strings.ToLower(host), Port: defaultConnectPort})
		unported++
		if hasPort {
			for i := len(endpoints) - unported; i < len(endpoints); i++ {
				endpoints[i].Port = port
			}
			unported = 0
		}
	}
	return endpoints, true
}

// splitHostPort splits `host[:port]`, with `[addr]:port` for IPv6.
func splitHostPort(s string) (host, port string, hasPort, ok bool) {
	host = s
	if strings.HasPrefix(s, "[") {
		end := strings.Index(s, "]")
		if end < 0 {
			return "", "", false, false
		}
		host = s[1:end]
		if tail := s[end+1:]; tail != "" {
			if !strings.HasPrefix(tail, ":") {
				return "", "", false, false
			}
			port, hasPort = tail[1:], true
		}
	} else if h, p, found := strings.Cut(s, ":"); found {
		host, port, hasPort = h, p, true
	}
	if host == "" || (hasPort && (port == "" || strings.Contains(port, ":"))) {
		return "", "", false, false
	}
	return host, port, hasPort, true
}

func parseTNSDescriptor(s string) (ConnectTarget, bool) {
	if block, closed := tnsBlock(s, 0); !closed || len(block) != len(s) {
		return ConnectTarget{}, false
	}
	if tnsSourceRoutePattern.MatchString(s) {
		return ConnectTarget{}, false
	}
	// A DESCRIPTION_LIST repeats CONNECT_DATA per description; the same identity
	// named twice is still one. A service and a SID, or two different values,
	// name two databases.
	var identities []string
	target := ConnectTarget{}
	for _, m := range tnsKeywordPattern.FindAllStringSubmatch(s, -1) {
		key := strings.ToUpper(m[1])
		if key != "SERVICE_NAME" && key != "SID" {
			continue
		}
		val := strings.TrimSpace(m[2])
		if val == "" {
			return ConnectTarget{}, false
		}
		if len(identities) == 0 {
			identities = append(identities, key+"="+val)
			if key == "SID" {
				target.SID = val
			} else {
				target.Service = val
			}
			continue
		}
		if !strings.EqualFold(identities[0], key+"="+val) {
			identities = append(identities, key+"="+val)
		}
	}
	if len(identities) != 1 {
		return ConnectTarget{}, false
	}
	var endpoints []Endpoint
	for _, loc := range tnsAddressPattern.FindAllStringIndex(s, -1) {
		block, closed := tnsBlock(s, loc[0])
		if !closed {
			return ConnectTarget{}, false
		}
		var hosts, ports []string
		for _, m := range tnsKeywordPattern.FindAllStringSubmatch(block, -1) {
			val := strings.TrimSpace(m[2])
			switch strings.ToUpper(m[1]) {
			case "HOST":
				if val == "" {
					return ConnectTarget{}, false
				}
				hosts = append(hosts, val)
			case "PORT":
				if val == "" {
					return ConnectTarget{}, false
				}
				ports = append(ports, val)
			default:
			}
		}
		if len(hosts) != 1 || len(ports) > 1 {
			return ConnectTarget{}, false
		}
		port := defaultConnectPort
		if len(ports) == 1 {
			port = ports[0]
		}
		endpoints = append(endpoints, Endpoint{Host: strings.ToLower(hosts[0]), Port: port})
	}
	if len(endpoints) == 0 {
		return ConnectTarget{}, false
	}
	target.Endpoints = endpoints
	return target, true
}

// tnsBlock returns the parenthesized block starting at s[start] and reports
// whether a matching closing parenthesis was found.
func tnsBlock(s string, start int) (string, bool) {
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return s[start : i+1], true
			}
		default:
		}
	}
	return "", false
}
