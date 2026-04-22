package parser

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseCreatePlacementPolicyStmt parses CREATE [OR REPLACE] PLACEMENT
// POLICY [IF NOT EXISTS] policy_name placement_option_list.
//
// Called from parseCreateDispatch after it has consumed CREATE and an
// optional OR REPLACE; p.cur is kwPLACEMENT on entry.
//
// Ref: TiDB v8.5.5 parser.y:15427-15436 (CreatePolicyStmt).
func (p *Parser) parseCreatePlacementPolicyStmt(start int, orReplace bool) (nodes.Node, error) {
	p.advance() // consume PLACEMENT
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &nodes.CreatePlacementPolicyStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Completion: after the IF NOT EXISTS clause, the user is naming a
	// fresh policy — no candidates to offer.
	p.checkCursor()
	if p.collectMode() {
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parsePlacementOptionList()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterPlacementPolicyStmt parses ALTER PLACEMENT POLICY
// [IF EXISTS] policy_name placement_option_list. Option list is shared
// with CREATE (DirectPlacementOption).
//
// Called from parseAlterDispatch; p.cur is kwPLACEMENT on entry.
//
// Ref: TiDB v8.5.5 parser.y:15438-15446.
func (p *Parser) parseAlterPlacementPolicyStmt(start int) (nodes.Node, error) {
	p.advance() // consume PLACEMENT
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &nodes.AlterPlacementPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("placement_policy_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parsePlacementOptionList()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseDropPlacementPolicyStmt parses DROP PLACEMENT POLICY
// [IF EXISTS] policy_name. The grammar does NOT allow a comma-list of
// policy names — unlike DROP TABLE.
//
// Called from parseDropDispatch; p.cur is kwPLACEMENT on entry.
//
// Ref: TiDB v8.5.5 parser.y:15389-15396.
func (p *Parser) parseDropPlacementPolicyStmt(start int) (nodes.Node, error) {
	p.advance() // consume PLACEMENT
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	stmt := &nodes.DropPlacementPolicyStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("placement_policy_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parsePlacementOptionList parses one-or-more DirectPlacementOption
// items, separated by commas OR whitespace (the TiDB grammar allows
// both and even a mix — parser.y:1999-2011). Returns when the next
// token is neither a placement option keyword nor a comma.
func (p *Parser) parsePlacementOptionList() ([]*nodes.PlacementPolicyOption, error) {
	var opts []*nodes.PlacementPolicyOption
	for {
		// Completion: offer option keywords at any position in the list.
		p.checkCursor()
		if p.collectMode() {
			for _, t := range placementOptionTokens {
				p.addTokenCandidate(t)
			}
			return opts, &ParseError{Message: "collecting"}
		}

		if !isPlacementOptionStart(p.cur.Type) {
			break
		}
		opt, err := p.parsePlacementOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		// Optional comma separator between options.
		if p.cur.Type == ',' {
			p.advance()
		}
	}
	return opts, nil
}

// placementOptionTokens is the closed set of keyword tokens that can
// start a DirectPlacementOption. Kept in one place so the completion
// walker and the option-list loop see the same set.
var placementOptionTokens = []int{
	kwPRIMARY_REGION, kwREGIONS, kwFOLLOWERS, kwVOTERS, kwLEARNERS,
	kwSCHEDULE, kwCONSTRAINTS,
	kwLEADER_CONSTRAINTS, kwFOLLOWER_CONSTRAINTS,
	kwVOTER_CONSTRAINTS, kwLEARNER_CONSTRAINTS,
	kwSURVIVAL_PREFERENCES,
}

func isPlacementOptionStart(t int) bool {
	for _, k := range placementOptionTokens {
		if t == k {
			return true
		}
	}
	return false
}

// parsePlacementOption parses a single DirectPlacementOption.
// `=` is optional for every option per EqOpt (parser.y:5450-5452).
func (p *Parser) parsePlacementOption() (*nodes.PlacementPolicyOption, error) {
	start := p.pos()
	tok := p.cur.Type
	name := ""
	isStringOpt := true
	isIntOpt := false

	switch tok {
	case kwPRIMARY_REGION:
		name = "PRIMARY_REGION"
	case kwREGIONS:
		name = "REGIONS"
	case kwSCHEDULE:
		name = "SCHEDULE"
	case kwCONSTRAINTS:
		name = "CONSTRAINTS"
	case kwLEADER_CONSTRAINTS:
		name = "LEADER_CONSTRAINTS"
	case kwFOLLOWER_CONSTRAINTS:
		name = "FOLLOWER_CONSTRAINTS"
	case kwVOTER_CONSTRAINTS:
		name = "VOTER_CONSTRAINTS"
	case kwLEARNER_CONSTRAINTS:
		name = "LEARNER_CONSTRAINTS"
	case kwSURVIVAL_PREFERENCES:
		name = "SURVIVAL_PREFERENCES"
	case kwFOLLOWERS:
		name = "FOLLOWERS"
		isStringOpt = false
		isIntOpt = true
	case kwVOTERS:
		name = "VOTERS"
		isStringOpt = false
		isIntOpt = true
	case kwLEARNERS:
		name = "LEARNERS"
		isStringOpt = false
		isIntOpt = true
	default:
		return nil, &ParseError{Message: fmt.Sprintf("unexpected token %d in PLACEMENT POLICY option", tok), Position: p.cur.Loc}
	}
	p.advance()
	p.match('=') // optional

	opt := &nodes.PlacementPolicyOption{
		Loc:  nodes.Loc{Start: start},
		Name: name,
	}

	if isStringOpt {
		if p.cur.Type != tokSCONST {
			return nil, &ParseError{Message: fmt.Sprintf("expected string literal for PLACEMENT POLICY option %s", name), Position: p.cur.Loc}
		}
		opt.Value = p.cur.Str
		p.advance()
	} else if isIntOpt {
		if p.cur.Type != tokICONST {
			return nil, &ParseError{Message: fmt.Sprintf("expected integer literal for PLACEMENT POLICY option %s", name), Position: p.cur.Loc}
		}
		if p.cur.Ival < 0 {
			return nil, &ParseError{Message: fmt.Sprintf("%s must be non-negative", name), Position: p.cur.Loc}
		}
		if name == "FOLLOWERS" && p.cur.Ival == 0 {
			// Matches upstream: parser.y:2025-2028 "FOLLOWERS must be positive".
			return nil, &ParseError{Message: "FOLLOWERS must be positive", Position: p.cur.Loc}
		}
		opt.IntValue = uint64(p.cur.Ival)
		opt.IsInt = true
		opt.Value = fmt.Sprintf("%d", p.cur.Ival)
		p.advance()
	}
	opt.Loc.End = p.pos()
	return opt, nil
}
