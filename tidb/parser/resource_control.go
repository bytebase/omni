package parser

import (
	"strconv"
	"strings"
	"time"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's resource control statements: CREATE, ALTER and DROP RESOURCE GROUP,
// SET RESOURCE GROUP, CALIBRATE RESOURCE and the two QUERY WATCH forms.
//
// None of this is MySQL's. MySQL 8.0 has CREATE/ALTER/DROP RESOURCE GROUP too,
// with a mandatory TYPE and the options VCPU, THREAD_PRIORITY and
// ENABLE/DISABLE; TiDB shares only the three statement names and answers 1064
// for every one of those options. The grammar below is TiDB v8.5.5's own, read
// from parser.y and confirmed form by form against a live v8.5.5 in
// resource_control_test.go.
//
// Ref: TiDB v8.5.5 parser.y:15733 (CREATE), 15743 (ALTER), 15753 (DROP),
// 16882 (CALIBRATE RESOURCE), 16980 and 17068 (QUERY WATCH)
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-create-resource-group/

// eqOpt consumes an optional '=', the grammar's EqOpt.
func (p *Parser) eqOpt() {
	p.match('=')
}

// parseResourceGroupName parses ResourceGroupName (parser.y:4495): an
// identifier, or the keyword DEFAULT, which names the group every session
// starts in.
func (p *Parser) parseResourceGroupName() (string, error) {
	if p.cur.Type == kwDEFAULT {
		tok := p.advance()
		return tok.Str, nil
	}
	// A user variable is not a group name here. The lexer scans @rg as one
	// tokIDENT, which parseIdent would take: v8.5.5 answers 1064 for
	// `ACTION SWITCH_GROUP(@rg)` and takes @rg only where the grammar names
	// UserVariable, which is after QUERY WATCH's RESOURCE GROUP.
	if p.atUserVariable() || p.cur.Type == tokAt2 {
		return "", &ParseError{
			Message:  "expected a resource group name",
			Position: p.cur.Loc,
		}
	}
	name, _, err := p.parseIdent()
	return name, err
}

// parseLengthNum parses the grammar's LengthNum: an unsigned integer literal,
// with no sign and no overflow past uint64.
func (p *Parser) parseLengthNum() (uint64, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	tok := p.cur
	val, err := strconv.ParseUint(tok.Str, 10, 64)
	if err != nil {
		return 0, &ParseError{
			Message:  "value out of range [0, 18446744073709551615]",
			Position: tok.Loc,
		}
	}
	p.advance()
	return val, nil
}

// parseStringLit parses a string literal, the grammar's stringLit.
func (p *Parser) parseStringLit() (string, int, error) {
	if p.cur.Type != tokSCONST {
		return "", 0, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return tok.Str, tok.Loc, nil
}

// goDuration reports whether Go's time.ParseDuration reads s, which is the
// check the grammar runs on EXEC_ELAPSED and on WATCH's DURATION
// (parser.y:1878, 1911). It is why `EXEC_ELAPSED '1d'` and `EXEC_ELAPSED '10'`
// are syntax errors on v8.5.5 while `'1h30m'` and `'-1s'` are not: Go's units
// stop at the hour and a number without one is not a duration.
func goDuration(s string) bool {
	_, err := time.ParseDuration(s)
	return err == nil
}

// tidbDuration reports whether TiDB's own duration parser reads s, which is
// the check the grammar runs on CALIBRATE RESOURCE's DURATION
// (parser.y:16938). That parser takes the units d, h and m only, reads a
// leading number only when a unit follows it, and special-cases "0"; the
// empty string parses because its loop never runs. So '20d' and '1d12h' are
// durations here where Go's parser has no 'd', and '20', '1w' and '-1m' are
// not.
//
// Ref: github.com/pingcap/tidb/pkg/parser/duration.ParseDuration
func tidbDuration(s string) bool {
	if s == "0" {
		return true
	}
	for len(s) > 0 {
		digits := 0
		for digits < len(s) && (s[digits] == '.' || (s[digits] >= '0' && s[digits] <= '9')) {
			digits++
		}
		// readFloat finds the number by looking for the first byte that is
		// neither a digit nor a dot, so a string of digits alone ends the
		// scan with nothing read and fails.
		if digits == 0 || digits == len(s) {
			return false
		}
		if _, err := strconv.ParseFloat(s[:digits], 64); err != nil {
			return false
		}
		switch s[digits] {
		case 'd', 'h', 'm':
		default:
			return false
		}
		s = s[digits+1:]
	}
	return true
}

// parseResourceGroupOptionList parses ResourceGroupOptionList
// (parser.y:1784). The options are separated by a comma or by nothing, at
// least one is required, and a repeat of one the list already carries is a
// syntax error (CheckAppend, parser.y:1791).
func (p *Parser) parseResourceGroupOptionList() ([]*nodes.ResourceGroupOption, error) {
	var options []*nodes.ResourceGroupOption
	for {
		option, err := p.parseResourceGroupOption()
		if err != nil {
			return nil, err
		}
		for _, seen := range options {
			if seen.Name == option.Name {
				return nil, &ParseError{
					Message:  "duplicate resource group option " + option.Name,
					Position: option.Loc.Start,
				}
			}
		}
		options = append(options, option)
		if p.cur.Type == ',' {
			p.advance() // consume ,
			continue
		}
		if !p.atResourceGroupOption() {
			break
		}
	}
	return options, nil
}

// atResourceGroupOption reports whether another resource group option starts
// here. The list takes a comma between options or nothing at all, so
// `ru_per_sec=1000 BURSTABLE` and `ru_per_sec=1000, BURSTABLE` are the same
// group.
func (p *Parser) atResourceGroupOption() bool {
	switch p.cur.Type {
	case kwRU_PER_SEC, kwPRIORITY, kwBURSTABLE, kwQUERY_LIMIT, kwBACKGROUND:
		return true
	}
	return false
}

// parseResourceGroupOption parses DirectResourceGroupOption (parser.y:1946).
func (p *Parser) parseResourceGroupOption() (*nodes.ResourceGroupOption, error) {
	start := p.pos()
	option := &nodes.ResourceGroupOption{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwRU_PER_SEC:
		option.Name = "RU_PER_SEC"
		p.advance()
		p.eqOpt()
		if p.cur.Type == kwUNLIMITED {
			p.advance()
			option.Unlimited = true
			break
		}
		val, err := p.parseLengthNum()
		if err != nil {
			return nil, err
		}
		option.IntValue = val

	case kwPRIORITY:
		option.Name = "PRIORITY"
		p.advance()
		p.eqOpt()
		switch p.cur.Type {
		case kwLOW, kwMEDIUM, kwHIGH:
			option.Value = strings.ToUpper(p.advance().Str)
		default:
			return nil, &ParseError{
				Message:  "expected LOW, MEDIUM or HIGH for PRIORITY",
				Position: p.cur.Loc,
			}
		}

	case kwBURSTABLE:
		// BURSTABLE on its own is the option set; v8.5.5 also takes a
		// boolean value, which pingcap's parser on master has replaced with
		// the words MODERATED, UNLIMITED and OFF (parser.y:1960-1974). Those
		// words are 1064 on v8.5.5, so BURSTABLE reads a number, TRUE or
		// FALSE and nothing else.
		option.Name = "BURSTABLE"
		p.advance()
		eq := p.cur.Type == '='
		p.eqOpt()
		switch p.cur.Type {
		case kwTRUE:
			p.advance()
			option.Burstable = true
		case kwFALSE:
			p.advance()
		case tokICONST:
			// The grammar reads this number as Int64Num, whose action answers 1105
			// for anything an int64 cannot hold.
			val, err := p.parseAdminInt64()
			if err != nil {
				return nil, err
			}
			option.Burstable = val != 0
		default:
			if eq {
				return nil, &ParseError{
					Message:  "expected TRUE, FALSE or a number for BURSTABLE",
					Position: p.cur.Loc,
				}
			}
			option.Burstable = true
		}

	case kwQUERY_LIMIT:
		option.Name = "QUERY_LIMIT"
		p.advance()
		p.eqOpt()
		// QUERY_LIMIT = NULL and QUERY_LIMIT = () both mean no rules.
		if p.cur.Type == kwNULL {
			p.advance()
			break
		}
		if _, ok := p.match('('); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			rules, err := p.parseRunawayOptionList()
			if err != nil {
				return nil, err
			}
			option.Runaway = rules
		}
		if _, ok := p.match(')'); !ok {
			return nil, p.syntaxErrorAtCur()
		}

	case kwBACKGROUND:
		option.Name = "BACKGROUND"
		p.advance()
		p.eqOpt()
		if p.cur.Type == kwNULL {
			p.advance()
			break
		}
		if _, ok := p.match('('); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			opts, err := p.parseBackgroundOptionList()
			if err != nil {
				return nil, err
			}
			option.Background = opts
		}
		if _, ok := p.match(')'); !ok {
			return nil, p.syntaxErrorAtCur()
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}
	option.Loc.End = p.prev.End
	return option, nil
}

// parseRunawayOptionList parses ResourceGroupRunawayOptionList
// (parser.y:1820), the rules inside QUERY_LIMIT. v8.5.5 refuses a second
// ACTION or WATCH and takes a second EXEC_ELAPSED, PROCESSED_KEYS or RU,
// keeping the last of them.
func (p *Parser) parseRunawayOptionList() ([]*nodes.ResourceGroupRunawayOption, error) {
	var options []*nodes.ResourceGroupRunawayOption
	for {
		option, err := p.parseRunawayOption()
		if err != nil {
			return nil, err
		}
		if option.Name == "ACTION" || option.Name == "WATCH" {
			for _, seen := range options {
				if seen.Name == option.Name {
					return nil, &ParseError{
						Message:  "duplicate QUERY_LIMIT option " + option.Name,
						Position: option.Loc.Start,
					}
				}
			}
		}
		options = append(options, option)
		if p.cur.Type == ',' {
			p.advance() // consume ,
			continue
		}
		if !p.atRunawayOption() {
			break
		}
	}
	return options, nil
}

// atRunawayOption reports whether another QUERY_LIMIT rule starts here.
func (p *Parser) atRunawayOption() bool {
	switch p.cur.Type {
	case kwEXEC_ELAPSED, kwPROCESSED_KEYS, kwRU, kwACTION, kwWATCH:
		return true
	}
	return false
}

// parseRunawayOption parses DirectResourceGroupRunawayOption
// (parser.y:1877).
func (p *Parser) parseRunawayOption() (*nodes.ResourceGroupRunawayOption, error) {
	start := p.pos()
	option := &nodes.ResourceGroupRunawayOption{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwEXEC_ELAPSED:
		option.Name = "EXEC_ELAPSED"
		p.advance()
		p.eqOpt()
		val, loc, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		if !goDuration(val) {
			return nil, &ParseError{
				Message:  "the EXEC_ELAPSED option is not a valid duration",
				Position: loc,
			}
		}
		option.Value = val

	case kwPROCESSED_KEYS, kwRU:
		// The grammar names intLit here, which the lexer scans unsigned and
		// without a sign of its own: `PROCESSED_KEYS -1` and `RU +1` are
		// 1064 on v8.5.5. The value is Int64Num's, so 2^63 answers 1105.
		option.Name = "PROCESSED_KEYS"
		if p.cur.Type == kwRU {
			option.Name = "RU"
		}
		p.advance()
		p.eqOpt()
		val, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		option.IntValue = uint64(val)

	case kwACTION:
		option.Name = "ACTION"
		p.advance()
		p.eqOpt()
		action, group, err := p.parseRunawayAction()
		if err != nil {
			return nil, err
		}
		option.Action, option.SwitchGroup = action, group

	case kwWATCH:
		option.Name = "WATCH"
		p.advance()
		p.eqOpt()
		watch, err := p.parseRunawayWatchType()
		if err != nil {
			return nil, err
		}
		option.Watch = watch
		// WatchDurationOption (parser.y:1936): absent, a duration string, or
		// the keyword UNLIMITED. The action lowercases the string before it
		// looks at it, so the string 'unlimited' in any case stands for the
		// keyword and skips the duration check, and `DURATION '1H'` and
		// `DURATION '1H30M'` are durations where `DURATION '1d'` is not --
		// Go's parser has no 'd' and reads only lower-case units.
		if p.cur.Type == kwDURATION {
			p.advance()
			p.eqOpt()
			if p.cur.Type == kwUNLIMITED {
				p.advance()
				option.DurationUnlimited = true
				break
			}
			val, loc, err := p.parseStringLit()
			if err != nil {
				return nil, err
			}
			lower := strings.ToLower(val)
			if lower == "unlimited" {
				// The two spellings land in the same place, as they do
				// upstream: the server keeps no duration for either.
				option.DurationUnlimited = true
				break
			}
			if lower != "" && !goDuration(lower) {
				return nil, &ParseError{
					Message:  "the WATCH DURATION option is not a valid duration",
					Position: loc,
				}
			}
			option.Duration = val
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}
	option.Loc.End = p.prev.End
	return option, nil
}

// parseRunawayAction parses ResourceGroupRunawayActionOption
// (parser.y:1856), returning the action and, for SWITCH_GROUP, the group it
// names.
func (p *Parser) parseRunawayAction() (string, string, error) {
	switch p.cur.Type {
	case kwDRYRUN, kwCOOLDOWN, kwKILL:
		return strings.ToUpper(p.advance().Str), "", nil
	case kwSWITCH_GROUP:
		p.advance()
		if _, ok := p.match('('); !ok {
			return "", "", p.syntaxErrorAtCur()
		}
		group, err := p.parseResourceGroupName()
		if err != nil {
			return "", "", err
		}
		if _, ok := p.match(')'); !ok {
			return "", "", p.syntaxErrorAtCur()
		}
		return "SWITCH_GROUP", group, nil
	}
	return "", "", &ParseError{
		Message:  "expected DRYRUN, COOLDOWN, KILL or SWITCH_GROUP for ACTION",
		Position: p.cur.Loc,
	}
}

// parseRunawayWatchType parses ResourceGroupRunawayWatchOption
// (parser.y:1842).
func (p *Parser) parseRunawayWatchType() (string, error) {
	switch p.cur.Type {
	case kwEXACT, kwSIMILAR, kwPLAN:
		return strings.ToUpper(p.advance().Str), nil
	}
	return "", &ParseError{
		Message:  "expected EXACT, SIMILAR or PLAN for WATCH",
		Position: p.cur.Loc,
	}
}

// parseBackgroundOptionList parses ResourceGroupBackgroundOptionList
// (parser.y:2000), the options inside BACKGROUND. A repeat is a syntax error
// (CheckBackgroundAppend, parser.y:2007).
func (p *Parser) parseBackgroundOptionList() ([]*nodes.ResourceGroupBackgroundOption, error) {
	var options []*nodes.ResourceGroupBackgroundOption
	for {
		start := p.pos()
		option := &nodes.ResourceGroupBackgroundOption{Loc: nodes.Loc{Start: start}}
		switch p.cur.Type {
		case kwTASK_TYPES:
			option.Name = "TASK_TYPES"
			p.advance()
			p.eqOpt()
			val, _, err := p.parseStringLit()
			if err != nil {
				return nil, err
			}
			option.Value = val
		case kwUTILIZATION_LIMIT:
			option.Name = "UTILIZATION_LIMIT"
			p.advance()
			p.eqOpt()
			val, err := p.parseLengthNum()
			if err != nil {
				return nil, err
			}
			option.IntValue = val
		default:
			return nil, p.syntaxErrorAtCur()
		}
		option.Loc.End = p.prev.End
		for _, seen := range options {
			if seen.Name == option.Name {
				return nil, &ParseError{
					Message:  "duplicate BACKGROUND option " + option.Name,
					Position: option.Loc.Start,
				}
			}
		}
		options = append(options, option)
		if p.cur.Type == ',' {
			p.advance() // consume ,
			continue
		}
		if p.cur.Type != kwTASK_TYPES && p.cur.Type != kwUTILIZATION_LIMIT {
			break
		}
	}
	return options, nil
}

// parseCreateResourceGroupStmt parses a CREATE RESOURCE GROUP statement.
// RESOURCE is already consumed by the CREATE dispatch.
//
// Ref: TiDB v8.5.5 parser.y:15733
//
//	CREATE RESOURCE GROUP [IF NOT EXISTS] group_name option [[,] option] ...
func (p *Parser) parseCreateResourceGroupStmt(start int) (*nodes.CreateResourceGroupStmt, error) {
	if _, ok := p.match(kwGROUP); !ok {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	stmt := &nodes.CreateResourceGroupStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance()
		if _, ok := p.match(kwNOT); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		if _, ok := p.match(kwEXISTS_KW); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.IfNotExists = true
	}
	name, err := p.parseResourceGroupName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	// The option list is not optional: v8.5.5 answers 1064 for a group with
	// no options at all.
	options, err := p.parseResourceGroupOptionList()
	if err != nil {
		return nil, err
	}
	stmt.Options = options
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterResourceGroupStmt parses an ALTER RESOURCE GROUP statement.
// RESOURCE is already consumed by the ALTER dispatch.
//
// Ref: TiDB v8.5.5 parser.y:15743
//
//	ALTER RESOURCE GROUP [IF EXISTS] group_name option [[,] option] ...
func (p *Parser) parseAlterResourceGroupStmt(start int) (*nodes.AlterResourceGroupStmt, error) {
	if _, ok := p.match(kwGROUP); !ok {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	stmt := &nodes.AlterResourceGroupStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance()
		if _, ok := p.match(kwEXISTS_KW); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.IfExists = true
	}
	name, err := p.parseResourceGroupName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	options, err := p.parseResourceGroupOptionList()
	if err != nil {
		return nil, err
	}
	stmt.Options = options
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropResourceGroupStmt parses a DROP RESOURCE GROUP statement.
// RESOURCE is already consumed by the DROP dispatch.
//
// Ref: TiDB v8.5.5 parser.y:15753
//
//	DROP RESOURCE GROUP [IF EXISTS] group_name
//
// One group, never a list: `DROP RESOURCE GROUP x, y` is 1064.
func (p *Parser) parseDropResourceGroupStmt(start int) (*nodes.DropResourceGroupStmt, error) {
	if _, ok := p.match(kwGROUP); !ok {
		return nil, &ParseError{
			Message:  "expected GROUP after RESOURCE",
			Position: p.cur.Loc,
		}
	}
	stmt := &nodes.DropResourceGroupStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance()
		if _, ok := p.match(kwEXISTS_KW); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.IfExists = true
	}
	name, err := p.parseResourceGroupName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCalibrateResourceStmt parses a CALIBRATE RESOURCE statement. CALIBRATE
// is the current token.
//
// Ref: TiDB v8.5.5 parser.y:16882
//
//	CALIBRATE RESOURCE [WORKLOAD workload | option [[,] option] ...]
//
//	workload: TPCC | OLTP_READ_WRITE | OLTP_READ_ONLY | OLTP_WRITE_ONLY | TPCH_10
//
//	option: START_TIME [=] expr
//	      | END_TIME [=] expr
//	      | DURATION [=] 'duration'
//	      | DURATION [=] INTERVAL expr unit
//
// The two forms do not mix: CalibrateOption takes a workload or a list of
// dynamic options, so `WORKLOAD tpcc START_TIME ...` is 1064. WORKLOAD takes
// no `=` either.
func (p *Parser) parseCalibrateResourceStmt() (*nodes.CalibrateResourceStmt, error) {
	start := p.pos()
	p.advance() // consume CALIBRATE
	if _, ok := p.match(kwRESOURCE); !ok {
		return nil, &ParseError{
			Message:  "expected RESOURCE after CALIBRATE",
			Position: p.cur.Loc,
		}
	}
	stmt := &nodes.CalibrateResourceStmt{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwWORKLOAD:
		p.advance()
		switch p.cur.Type {
		case kwTPCC, kwOLTP_READ_WRITE, kwOLTP_READ_ONLY, kwOLTP_WRITE_ONLY, kwTPCH_10:
			stmt.Workload = strings.ToUpper(p.advance().Str)
		default:
			return nil, &ParseError{
				Message:  "expected TPCC, OLTP_READ_WRITE, OLTP_READ_ONLY, OLTP_WRITE_ONLY or TPCH_10 for WORKLOAD",
				Position: p.cur.Loc,
			}
		}
	case kwSTART_TIME, kwEND_TIME, kwDURATION:
		options, err := p.parseCalibrateOptionList()
		if err != nil {
			return nil, err
		}
		stmt.Options = options
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCalibrateOptionList parses DynamicCalibrateOptionList
// (parser.y:16905).
//
// The duplicate check is upstream's, which compares a new option against the
// first two of the list and no further (parser.y:16911). Three option kinds
// exist, so the only list it lets a repeat into is one that already holds
// three: `START_TIME a END_TIME b DURATION c DURATION d` parses on v8.5.5
// while `DURATION c START_TIME a DURATION d` does not.
func (p *Parser) parseCalibrateOptionList() ([]*nodes.CalibrateResourceOption, error) {
	var options []*nodes.CalibrateResourceOption
	for {
		option, err := p.parseCalibrateOption()
		if err != nil {
			return nil, err
		}
		for i, seen := range options {
			if i > 1 {
				break
			}
			if seen.Name == option.Name {
				return nil, &ParseError{
					Message:  "duplicate CALIBRATE RESOURCE option " + option.Name,
					Position: option.Loc.Start,
				}
			}
		}
		options = append(options, option)
		if p.cur.Type == ',' {
			p.advance() // consume ,
			continue
		}
		if p.cur.Type != kwSTART_TIME && p.cur.Type != kwEND_TIME && p.cur.Type != kwDURATION {
			break
		}
	}
	return options, nil
}

// parseCalibrateOption parses DynamicCalibrateResourceOption
// (parser.y:16929).
func (p *Parser) parseCalibrateOption() (*nodes.CalibrateResourceOption, error) {
	start := p.pos()
	option := &nodes.CalibrateResourceOption{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwSTART_TIME, kwEND_TIME:
		option.Name = "START_TIME"
		if p.cur.Type == kwEND_TIME {
			option.Name = "END_TIME"
		}
		p.advance()
		p.eqOpt()
		// A whole Expression, not a literal: v8.5.5 takes
		// `START_TIME now() - interval 20 minute` and `START_TIME a IS NULL`.
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		option.Ts = expr

	case kwDURATION:
		option.Name = "DURATION"
		p.advance()
		p.eqOpt()
		if p.cur.Type == kwINTERVAL {
			p.advance()
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			// TimeUnit is a keyword, not an identifier (parser.y:8511), so
			// a quoted one is not a unit: v8.5.5 answers 1064 for
			// `DURATION INTERVAL 15 ` + "`MINUTE`" + `. parseTimeUnitKeyword
			// turns away anything the lexer did not give a keyword token,
			// which a backtick-quoted word is not.
			unit, err := p.parseTimeUnitKeyword(false)
			if err != nil {
				return nil, err
			}
			option.Ts, option.Unit = expr, unit.Keyword
			break
		}
		val, loc, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		if !tidbDuration(val) {
			return nil, &ParseError{
				Message:  "the DURATION option is not a valid duration",
				Position: loc,
			}
		}
		option.Value = val

	default:
		return nil, p.syntaxErrorAtCur()
	}
	option.Loc.End = p.prev.End
	return option, nil
}

// parseQueryWatchStmt parses the two QUERY WATCH statements. QUERY is the
// current token and the caller has checked that WATCH follows it.
//
// Ref: TiDB v8.5.5 parser.y:16980 (ADD), 17068 (REMOVE)
//
//	QUERY WATCH ADD option [[,] option] ...
//	QUERY WATCH REMOVE watch_id
func (p *Parser) parseQueryWatchStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume QUERY
	p.advance() // consume WATCH
	switch p.cur.Type {
	case kwADD:
		p.advance()
		options, err := p.parseQueryWatchOptionList()
		if err != nil {
			return nil, err
		}
		return &nodes.AddQueryWatchStmt{
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
			Options: options,
		}, nil
	case kwREMOVE:
		p.advance()
		// A watch id and nothing else. The RESOURCE GROUP forms pingcap's
		// parser on master reads (parser.y:17075, 17081) are 1064 on v8.5.5,
		// and so is a negative id, and the id is Int64Num's, so 2^63 answers 1105.
		id, err := p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
		return &nodes.DropQueryWatchStmt{
			Loc: nodes.Loc{Start: start, End: p.prev.End},
			ID:  uint64(id),
		}, nil
	}
	return nil, &ParseError{
		Message:  "expected ADD or REMOVE after QUERY WATCH",
		Position: p.cur.Loc,
	}
}

// parseQueryWatchOptionList parses QueryWatchOptionList (parser.y:16988). One
// option of each kind at most, and the three text kinds share one slot
// (CheckQueryWatchAppend, parser.y:16995).
func (p *Parser) parseQueryWatchOptionList() ([]*nodes.QueryWatchOption, error) {
	var options []*nodes.QueryWatchOption
	for {
		option, err := p.parseQueryWatchOption()
		if err != nil {
			return nil, err
		}
		for _, seen := range options {
			if queryWatchSlot(seen.Name) == queryWatchSlot(option.Name) {
				return nil, &ParseError{
					Message:  "duplicate QUERY WATCH option " + option.Name,
					Position: option.Loc.Start,
				}
			}
		}
		options = append(options, option)
		if p.cur.Type == ',' {
			p.advance() // consume ,
			continue
		}
		if !p.atQueryWatchOption() {
			break
		}
	}
	return options, nil
}

// queryWatchSlot returns the slot an option name occupies. SQL DIGEST, PLAN
// DIGEST and SQL TEXT are three spellings of one option, so a statement
// carries one of them.
func queryWatchSlot(name string) string {
	switch name {
	case "SQL_DIGEST", "PLAN_DIGEST", "SQL_TEXT":
		return "TEXT"
	}
	return name
}

// atUserVariable reports whether the current token is a user variable, @name,
// which the lexer scans as one tokIDENT whose Str keeps the @. A system
// variable, @@name, is a different token and is not one.
//
// A backtick-quoted identifier reaches the parser as the same kind of token
// with the same kind of Str, so the @ in its text says nothing: a quoted
// @rg is an ordinary name, and v8.5.5 runs DROP RESOURCE GROUP and ACTION
// SWITCH_GROUP() with one. Only the quoting tells the two apart, which is
// what backtickQuotedAt reads back from the input.
func (p *Parser) atUserVariable() bool {
	return p.cur.Type == tokIDENT &&
		strings.HasPrefix(p.cur.Str, "@") &&
		!p.backtickQuotedAt(p.cur.Loc)
}

// atQueryWatchOption reports whether another QUERY WATCH option starts here.
func (p *Parser) atQueryWatchOption() bool {
	switch p.cur.Type {
	case kwRESOURCE, kwACTION, kwSQL, kwPLAN:
		return true
	}
	return false
}

// parseWatchPattern parses the SimpleExpr every QUERY WATCH pattern takes
// (QueryWatchTextOption, parser.y:17044).
//
// SimpleExpr carries two suffixes of its own that the precedence loop in
// parseExprPrec applies and parsePrefixExpr does not:
//
//   - JSON extraction, `SimpleIdent -> 'path'` and `SimpleIdent ->> 'path'`
//     (parser.y:8272, 8277). Its operands are narrower than the rest of
//     SimpleExpr: a column reference on the left, a string literal on the
//     right. That is why `cfg->>'$.d'` and `t.cfg->>'$.d'` parse on v8.5.5
//     while `'a'->>'$.d'`, `(cfg)->>'$.d'`, `now()->>'$.d'`, `cfg->>1` and
//     a chain of two are each 1064: the production's result is no longer a
//     SimpleIdent, so nothing can follow it.
//   - COLLATE (parser.y:8076), which takes a whole SimpleExpr and so repeats,
//     and may follow the JSON suffix where it may not precede it.
//
// The result is still a SimpleExpr and not a whole Expression, so
// `TO 'a' AND 1` stays 1064.
func (p *Parser) parseWatchPattern() (nodes.ExprNode, error) {
	start := p.pos()
	expr, err := p.parsePrefixExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == tokJsonExtract || p.cur.Type == tokJsonUnquote {
		if _, ok := expr.(*nodes.ColumnRef); !ok {
			return nil, &ParseError{
				Message:  "expected a column reference before the JSON path operator",
				Position: p.cur.Loc,
			}
		}
		op := nodes.BinOpJsonExtract
		if p.cur.Type == tokJsonUnquote {
			op = nodes.BinOpJsonUnquote
		}
		p.advance() // consume -> or ->>
		path, loc, err := p.parseStringLit()
		if err != nil {
			return nil, err
		}
		expr = &nodes.BinaryExpr{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Op:    op,
			Left:  expr,
			Right: &nodes.StringLit{Loc: nodes.Loc{Start: loc, End: p.prev.End}, Value: path},
		}
	}
	for p.cur.Type == kwCOLLATE {
		p.advance() // consume COLLATE
		collation, err := p.parseCollationName()
		if err != nil {
			return nil, err
		}
		expr = &nodes.CollateExpr{
			Loc:       nodes.Loc{Start: start, End: p.prev.End},
			Expr:      expr,
			Collation: collation,
		}
	}
	return expr, nil
}

// parseQueryWatchOption parses QueryWatchOption (parser.y:17010) and
// QueryWatchTextOption (parser.y:17044).
func (p *Parser) parseQueryWatchOption() (*nodes.QueryWatchOption, error) {
	start := p.pos()
	option := &nodes.QueryWatchOption{Loc: nodes.Loc{Start: start}}
	switch p.cur.Type {
	case kwRESOURCE:
		option.Name = "RESOURCE_GROUP"
		p.advance()
		if _, ok := p.match(kwGROUP); !ok {
			return nil, &ParseError{
				Message:  "expected GROUP after RESOURCE",
				Position: p.cur.Loc,
			}
		}
		// The grammar names UserVariable here (parser.y:17020), so @rg is a
		// group name and @@rg is not: v8.5.5 answers 1064 for
		// `RESOURCE GROUP @@global.rg`.
		if p.atUserVariable() {
			expr, err := p.parsePrefixExpr()
			if err != nil {
				return nil, err
			}
			option.ResourceGroupExpr = expr
			break
		}
		name, err := p.parseResourceGroupName()
		if err != nil {
			return nil, err
		}
		option.ResourceGroup = name

	case kwACTION:
		option.Name = "ACTION"
		p.advance()
		p.eqOpt()
		action, group, err := p.parseRunawayAction()
		if err != nil {
			return nil, err
		}
		option.Action, option.SwitchGroup = action, group

	case kwSQL:
		p.advance()
		switch p.cur.Type {
		case kwDIGEST:
			option.Name = "SQL_DIGEST"
			p.advance()
		case kwTEXT:
			option.Name = "SQL_TEXT"
			p.advance()
			watch, err := p.parseRunawayWatchType()
			if err != nil {
				return nil, err
			}
			option.Watch = watch
			if _, ok := p.match(kwTO); !ok {
				return nil, p.syntaxErrorAtCur()
			}
		default:
			return nil, &ParseError{
				Message:  "expected DIGEST or TEXT after SQL",
				Position: p.cur.Loc,
			}
		}
		expr, err := p.parseWatchPattern()
		if err != nil {
			return nil, err
		}
		option.Pattern = expr

	case kwPLAN:
		option.Name = "PLAN_DIGEST"
		p.advance()
		if _, ok := p.match(kwDIGEST); !ok {
			return nil, &ParseError{
				Message:  "expected DIGEST after PLAN",
				Position: p.cur.Loc,
			}
		}
		expr, err := p.parseWatchPattern()
		if err != nil {
			return nil, err
		}
		option.Pattern = expr

	default:
		return nil, p.syntaxErrorAtCur()
	}
	option.Loc.End = p.prev.End
	return option, nil
}
