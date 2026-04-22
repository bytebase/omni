package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// forXmlModes is SqlScriptDOM's ForXmlMode enum (RAW/AUTO/EXPLICIT/PATH).
// RAW and PATH are also registered keyword tokens.
var forXmlModes = newOptionSet(kwRAW, kwPATH).withIdents("RAW", "AUTO", "EXPLICIT", "PATH")

// forJsonModes is SqlScriptDOM's ForJsonMode enum (AUTO/PATH).
var forJsonModes = newOptionSet(kwPATH).withIdents("AUTO", "PATH")

// commaListFlags controls leniency of parseCommaList.
type commaListFlags uint8

const (
	// commaListStrict rejects both an empty list and a trailing comma.
	// This matches SqlScriptDOM behavior for the vast majority of T-SQL
	// comma-separated constructs (column lists, IN operator, VALUES row,
	// ORDER/GROUP BY, CTE list, PIVOT IN list, etc.).
	commaListStrict commaListFlags = 0

	// commaListAllowEmpty permits zero items. Used for call-site forms that
	// genuinely may have no items:
	//   - function / TVF argument lists: `GETDATE()`, `f()`
	//   - a single empty grouping set inside GROUPING SETS: `GROUPING SETS (..., (), ...)`
	// A trailing comma is still rejected.
	commaListAllowEmpty commaListFlags = 1 << 0

	// commaListAllowTrail permits a trailing comma before the terminator.
	// Reserved for constructs where SqlScriptDOM explicitly tolerates it; at
	// present only the top-level column list in CREATE TABLE qualifies.
	// An empty list is still rejected.
	commaListAllowTrail commaListFlags = 1 << 1
)

// parseCommaList drives a comma-separated list of items terminated by `end`.
// The opening delimiter (e.g. '(') must have been consumed by the caller;
// this helper does NOT consume the terminator either — the caller is
// responsible for `p.expect(end)` after we return.
//
// parseItem is invoked at the head of every item slot. Contract:
//   - If the current position starts a valid item, parseItem MUST advance past
//     that item and return it non-nil with err=nil.
//   - If the current position does NOT start a valid item, parseItem MUST
//     return an error (or the completion sentinel errCollecting).
//
// parseCommaList enforces this via a "must advance" guard: if parseItem
// returns (nil, nil) OR returns a non-nil node without having advanced the
// cursor, the helper treats it as an unexpected token. This catches
// item-level leniency bugs uniformly — many hand-written parseX helpers in
// this codebase follow a "try-parse returns nil on no match" convention that
// is unsafe in a strict-list context, and the guard keeps parseCommaList
// from silently appending empty/phantom entries.
//
// Leniency at the list boundary is governed by `flags`; see commaList*
// constants above.
//
// By design, parseCommaList is the single chokepoint for list shape
// validation in the mssql parser, so that strictness is uniform and
// regressions show up in one place.
func (p *Parser) parseCommaList(
	end int,
	flags commaListFlags,
	parseItem func() (nodes.Node, error),
) ([]nodes.Node, error) {
	if p.cur.Type == end {
		if flags&commaListAllowEmpty == 0 {
			return nil, p.unexpectedToken()
		}
		return nil, nil
	}
	var items []nodes.Node
	for {
		before := p.cur.Loc
		item, err := parseItem()
		if err != nil {
			return nil, err
		}
		// Enforce the parseItem contract: a successful call must both
		// produce a non-nil item AND advance the cursor. Anything else is
		// an item-level leniency bug on the caller side; surface it as
		// a plain syntax error at the current position.
		if item == nil || p.cur.Loc == before {
			return nil, p.unexpectedToken()
		}
		items = append(items, item)
		if _, ok := p.match(','); !ok {
			break
		}
		// Reject a trailing comma: after consuming ',' we must see another item
		// before the terminator. parseItem is allowed to run at EOF so it can
		// emit completion candidates via errCollecting; it must otherwise
		// return an error if there is no item to parse.
		if p.cur.Type == end {
			if flags&commaListAllowTrail == 0 {
				return nil, p.unexpectedToken()
			}
			break
		}
	}
	return items, nil
}
