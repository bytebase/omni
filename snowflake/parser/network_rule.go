package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Network rule DDL — CREATE / ALTER / DROP NETWORK RULE (gap-network-rule)
// ---------------------------------------------------------------------------
//
// A network rule carries a small, version-growing set of `KEY = value`
// properties (TYPE, VALUE_LIST, MODE, COMMENT, ...). The TYPE and MODE
// vocabularies grow with the service (IPV4 / IPV6 / AWSVPCEID / AZURELINKID /
// GCPPSCID / HOST_PORT / PRIVATE_HOST_PORT / COMPUTE_POOL / ...; INGRESS /
// EGRESS / INTERNAL_STAGE / SNOWFLAKE_MANAGED_STORAGE_VOLUME / ...). Rather
// than enumerate them, every property is parsed as an open-ended `KEY =
// <value>` pair (ast.CopyOption), reusing the COPY (T5.2) machinery
// (parseCopyOption / startsCopyOption) exactly as STAGE (T4.1), FILE FORMAT
// (T4.2) and WAREHOUSE (gap-warehouse) do. The parenthesized
// VALUE_LIST = ( '<str>' [ , ... ] ) string list is already a native
// CopyOption value shape (List). Property order is free; the catalog/semantic
// layer, not the parser, validates that a property is real.
//
// NETWORK RULE is distinguished from the pre-existing NETWORK POLICY object at
// dispatch: NETWORK is a reserved keyword but RULE is not, so `NETWORK RULE`
// lexes as kwNETWORK followed by a "RULE" identifier (curIsWord), whereas
// `NETWORK POLICY` lexes as kwNETWORK followed by kwPOLICY.

// networkRuleOptionCap bounds an open-ended network-rule property run. Real
// statements carry a handful of properties; a far larger run signals a runaway
// parse and aborts loudly rather than spinning.
const networkRuleOptionCap = 1024

// ---------------------------------------------------------------------------
// CREATE NETWORK RULE
// ---------------------------------------------------------------------------

// parseCreateNetworkRuleStmt parses the body of a
//
//	CREATE [ OR REPLACE ] NETWORK RULE [ IF NOT EXISTS ] <name>
//	  <prop> [ <prop> ... ]
//
// statement, where each <prop> is an open-ended `KEY = value` pair (TYPE,
// VALUE_LIST, MODE, COMMENT, ...). The CREATE keyword and the optional OR
// REPLACE modifier have already been consumed by parseCreateStmt; start is the
// Loc of the CREATE token, and cur is the NETWORK keyword.
func (p *Parser) parseCreateNetworkRuleStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume NETWORK
	// RULE is not a reserved keyword; the dispatcher confirmed it via curIsWord.
	if !p.curIsWord("RULE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume RULE

	stmt := &ast.CreateNetworkRuleStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS (lenient; the official grammar omits it, but it is accepted
	// here to mirror the other CREATE-object parsers).
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance() // consume IF
		p.advance() // consume NOT
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Network rule name (optionally qualified db.schema.name).
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Open-ended `KEY = value` property run. The loop ends at a statement
	// boundary or any token that does not begin a property. The loop-guard
	// aborts if an iteration fails to consume any token.
	opts, err := p.parseNetworkRuleOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER NETWORK RULE
// ---------------------------------------------------------------------------

// parseAlterNetworkRuleStmt parses
//
//	ALTER NETWORK RULE [ IF EXISTS ] <name> { SET <prop>... | UNSET <key>,... }
//
// The ALTER keyword has already been consumed; cur is the NETWORK keyword.
// SET carries an open-ended `KEY = value` property run (VALUE_LIST, COMMENT,
// ...); UNSET carries a comma-separated key-name list (VALUE_LIST, COMMENT,
// ...). The surface is kept lenient and open-ended, mirroring ALTER WAREHOUSE.
func (p *Parser) parseAlterNetworkRuleStmt() (ast.Node, error) {
	netTok := p.advance() // consume NETWORK (anchors Loc.Start at the object keyword)
	if !p.curIsWord("RULE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume RULE

	stmt := &ast.AlterNetworkRuleStmt{Loc: ast.Loc{Start: netTok.Loc.Start}}

	// Optional IF EXISTS.
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	// Network rule name.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch.
	switch p.cur.Type {
	case kwSET:
		p.advance() // consume SET
		opts, err := p.parseNetworkRuleOptions()
		if err != nil {
			return nil, err
		}
		if len(opts) == 0 {
			// SET with nothing settable is a syntax error.
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterNetworkRuleSet
		stmt.Options = opts

	case kwUNSET:
		p.advance() // consume UNSET
		keys, err := p.parseNetworkRuleUnsetKeys()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterNetworkRuleUnset
		stmt.UnsetKeys = keys

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared network-rule helpers
// ---------------------------------------------------------------------------

// parseNetworkRuleOptions parses a run of zero or more network-rule properties,
// each an open-ended `KEY = <value>` pair (reusing parseCopyOption). The run
// continues while the current token begins another property and terminates at a
// statement boundary or any token that does not begin a property. A loop-guard
// aborts if an iteration consumes no tokens.
func (p *Parser) parseNetworkRuleOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for i := 0; p.startsCopyOption(); i++ {
		if i >= networkRuleOptionCap {
			return nil, p.syntaxErrorAtCur()
		}
		before := p.cur.Loc.Start
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if p.cur.Loc.Start == before && p.cur.Type != tokEOF {
			// No forward progress on a non-EOF token: abort rather than spin.
			return nil, p.syntaxErrorAtCur()
		}
	}
	return opts, nil
}

// parseNetworkRuleUnsetKeys parses an UNSET <key> [ , ... ] property-name list.
// Each key is a single name word (identifier or keyword), uppercased; a comma
// separates keys. Consumes at least one key.
func (p *Parser) parseNetworkRuleUnsetKeys() ([]string, error) {
	var keys []string
	for {
		if !p.isOptionWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		keys = append(keys, strings.ToUpper(p.advance().Str))
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return keys, nil
}
