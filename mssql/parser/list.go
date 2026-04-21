package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

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
// parseItem is invoked at the head of every item slot. It must parse
// exactly one item and return it, or return an error if the current token
// does not start a valid item.
//
// Leniency is governed by `flags`; see commaList* constants above.
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
		item, err := parseItem()
		if err != nil {
			return nil, err
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
