package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// profileParamKind is how a profile parameter takes its value.
type profileParamKind int

const (
	// profileInteger takes { integer | UNLIMITED | DEFAULT }.
	profileInteger profileParamKind = iota
	// profileSize takes { integer [ K | M | G ] | UNLIMITED | DEFAULT }.
	profileSize
	// profileExpr takes { expr | UNLIMITED | DEFAULT }.
	profileExpr
	// profileFunction takes { function_name | NULL | DEFAULT }.
	profileFunction
)

// profileParams are the resource and password parameters of CREATE and
// ALTER PROFILE. Oracle Free 23.26.3 reads INACTIVE_ACCOUNT_TIME as an
// expression, as it does the other day counts, though the page shows an
// integer.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PROFILE.html
var profileParams = map[string]profileParamKind{
	"SESSIONS_PER_USER":         profileInteger,
	"CPU_PER_SESSION":           profileInteger,
	"CPU_PER_CALL":              profileInteger,
	"CONNECT_TIME":              profileInteger,
	"IDLE_TIME":                 profileInteger,
	"LOGICAL_READS_PER_SESSION": profileInteger,
	"LOGICAL_READS_PER_CALL":    profileInteger,
	"PRIVATE_SGA":               profileSize,
	"COMPOSITE_LIMIT":           profileInteger,
	"FAILED_LOGIN_ATTEMPTS":     profileInteger,
	"PASSWORD_REUSE_MAX":        profileInteger,
	"PASSWORD_LIFE_TIME":        profileExpr,
	"PASSWORD_REUSE_TIME":       profileExpr,
	"PASSWORD_LOCK_TIME":        profileExpr,
	"PASSWORD_GRACE_TIME":       profileExpr,
	"PASSWORD_ROLLOVER_TIME":    profileExpr,
	"INACTIVE_ACCOUNT_TIME":     profileExpr,
	"PASSWORD_VERIFY_FUNCTION":  profileFunction,
}

// isProfileParamStart reports whether the current token names a profile
// parameter.
func (p *Parser) isProfileParamStart() bool {
	if p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return false
	}
	_, ok := profileParams[p.cur.Str]
	return ok
}

// parseProfileLimits parses what follows the profile's name in CREATE and
// ALTER PROFILE:
//
//	LIMIT { resource_parameters | password_parameters }
//	      [ { resource_parameters | password_parameters } ]...
//	[ CONTAINER = { ALL | CURRENT } ]
//
// Oracle takes one LIMIT: a second one is not a parameter (ORA-02376).
func (p *Parser) parseProfileLimits(limits *[]*nodes.ProfileLimit, container **bool) error {
	if p.cur.Type != kwLIMIT {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	seen := map[string]bool{}
	for {
		if !p.isProfileParamStart() {
			return p.syntaxErrorAtCur()
		}
		lim, err := p.parseProfileParam(seen)
		if err != nil {
			return err
		}
		*limits = append(*limits, lim)
		seen[lim.Name] = true
		if !p.isProfileParamStart() {
			break
		}
	}
	if p.isKeywordStr("CONTAINER") {
		v, err := p.parseSecurityContainer()
		if err != nil {
			return err
		}
		*container = v
	}
	return nil
}

// parseProfileParam parses one profile parameter and its value. seen holds
// the parameters read before it.
func (p *Parser) parseProfileParam(seen map[string]bool) (*nodes.ProfileLimit, error) {
	start := p.pos()
	lim := &nodes.ProfileLimit{Name: p.cur.Str, Loc: nodes.Loc{Start: start}}
	kind := profileParams[lim.Name]
	p.advance()
	valTok := p.cur
	switch {
	case p.cur.Type == tokEOF || p.cur.Type == ';':
		// compat.tsv records a parameter named a second time with no value
		// at the end of the text as accepted, so omni takes it.
		if !seen[lim.Name] {
			return nil, p.syntaxErrorAtCur()
		}
		lim.Loc.End = p.prev.End
		return lim, nil
	case p.isKeywordStr("UNLIMITED") && kind != profileFunction, p.cur.Type == kwDEFAULT:
		lim.Value = p.advance().Str
	case kind == profileFunction:
		if p.cur.Type == kwNULL {
			lim.Value = p.advance().Str
			break
		}
		if p.cur.Type != tokQIDENT && (!p.isIdentLike() || isOracleSQLReservedKeyword(p.cur)) {
			return nil, p.syntaxErrorAtCur()
		}
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		lim.Value = name.Name
		if name.Schema != "" {
			lim.Value = name.Schema + "." + name.Name
		}
	case kind == profileExpr:
		exprStart := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		lim.Value = p.source[exprStart:p.prev.End]
	case p.cur.Type == tokICONST:
		valueTok := p.advance()
		lim.Value = valueTok.Str
		if kind == profileSize {
			switch {
			case p.isKeywordStr("K") || p.isKeywordStr("M") || p.isKeywordStr("G"):
				lim.Value += p.advance().Str
			case p.isKeywordStr("T") || p.isKeywordStr("P") || p.isKeywordStr("E"):
				// Oracle reads T, P and E, units the page lists, and
				// refuses the size at its integer (ORA-02377).
				return nil, p.syntaxErrorAtTok(valueTok)
			}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if seen[lim.Name] {
		return nil, p.syntaxErrorAtTok(valTok)
	}
	lim.Loc.End = p.prev.End
	return lim, nil
}

// parseSecurityContainer parses CONTAINER = { ALL | CURRENT } and returns
// true for ALL. The current token is CONTAINER.
func (p *Parser) parseSecurityContainer() (*bool, error) {
	p.advance() // CONTAINER
	if p.cur.Type != '=' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var all bool
	switch {
	case p.cur.Type == kwALL:
		all = true
	case p.isKeywordStr("CURRENT"):
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return &all, nil
}
